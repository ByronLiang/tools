package wf

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var defaultWaitTime = 60 * time.Second

// Watcher describes a process that watches files for changes.
type Watcher struct {
	// mu protects the following.
	mu       sync.RWMutex
	Event    chan Event
	Error    chan error
	closed   chan struct{}
	done     chan struct{}
	files    map[string]os.FileInfo // map of files.
	waitTime time.Duration
}

// New creates a new Watcher.
func NewPolling(waitTime time.Duration) *Watcher {
	// less 1 nanosecond use default duration: 60s
	if waitTime < time.Nanosecond {
		waitTime = defaultWaitTime
	}
	watch := &Watcher{
		Event:    make(chan Event),
		Error:    make(chan error),
		closed:   make(chan struct{}),
		done:     make(chan struct{}),
		files:    make(map[string]os.FileInfo),
		waitTime: waitTime,
	}
	go watch.readEvents()
	return watch
}

func (w *Watcher) readEvents() {
	defer close(w.done)
	timer := time.NewTimer(w.waitTime)
	if !timer.Stop() {
		<-timer.C
	}
	for {
		if w.isClosed() {
			return
		}
		timer.Reset(w.waitTime)
		select {
		case <-timer.C:
			fileList := w.retrieveFileList()
			w.pollEvent(fileList)
			w.updateFile(fileList)
		case <-w.closed:
			return
		}
	}
}

func (w *Watcher) isClosed() bool {
	select {
	case <-w.closed:
		return true
	default:
		return false
	}
}

func (w *Watcher) Close() {
	w.mu.Lock()
	if w.isClosed() {
		w.mu.Unlock()
		return
	}
	close(w.closed)
	w.mu.Unlock()
	// readEvent goroutine close
	<-w.done
	return
}

func (w *Watcher) Add(name string) error {
	name = filepath.Clean(name)
	if w.isClosed() {
		return errors.New("pollingWatcher instance already closed")
	}
	fileList, err := w.list(name)
	if err != nil {
		return err
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	for k, v := range fileList {
		w.files[k] = v
	}
	return nil
}

// 检测目录、文件名, 生成文件对象映射
func (w *Watcher) list(name string) (map[string]os.FileInfo, error) {
	fileList := make(map[string]os.FileInfo)

	// Make sure name exists.
	stat, err := os.Stat(name)
	if err != nil {
		return nil, err
	}

	fileList[name] = stat

	// If it's not a directory, just return.
	if !stat.IsDir() {
		return fileList, nil
	}

	// It's a directory.
	fInfoList, err := ioutil.ReadDir(name)
	if err != nil {
		return nil, err
	}
	// Add all of the files in the directory
	for _, fInfo := range fInfoList {
		path := filepath.Join(name, fInfo.Name())
		fileList[path] = fInfo
	}
	return fileList, nil
}

// Remove removes either a single file or directory from the file's list.
func (w *Watcher) Remove(name string) (err error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	name = filepath.Clean(name)

	// If name is a single file, remove it and return.
	info, found := w.files[name]
	if !found {
		return nil // Doesn't exist, just return.
	}
	if !info.IsDir() {
		delete(w.files, name)
		return nil
	}

	// Delete the actual directory from w.files
	delete(w.files, name)

	// If it's a directory, delete all of it's contents from w.files.
	for path := range w.files {
		if filepath.Dir(path) == name {
			delete(w.files, path)
		}
	}
	return nil
}

func (w *Watcher) updateFile(files map[string]os.FileInfo) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.files = files
	return
}

func (w *Watcher) retrieveFileList() map[string]os.FileInfo {
	w.mu.RLock()
	defer w.mu.RUnlock()

	fileList := make(map[string]os.FileInfo)

	var list map[string]os.FileInfo
	var err error

	for name := range w.files {
		list, err = w.list(name)
		if err != nil {
			// 针对非文件不存在错误, 向外传递
			if !os.IsNotExist(err) {
				w.Error <- err
			}
			continue
		}
		// Add the file's to the file list.
		for k, v := range list {
			fileList[k] = v
		}
	}

	return fileList
}

func (w *Watcher) pollEvent(files map[string]os.FileInfo) {
	w.mu.RLock()
	defer w.mu.RUnlock()
	// 检测被删除文件
	for path, info := range w.files {
		if _, found := files[path]; !found {
			select {
			case w.Event <- Event{Remove, path, info}:
			}
		}
	}

	// 检测新加文件与修改过文件
	for path, info := range files {
		oldInfo, found := w.files[path]
		if !found {
			// 新加文件
			select {
			case w.Event <- Event{Create, path, info}:
			}
			continue
		}
		if oldInfo.ModTime() != info.ModTime() {
			select {
			case w.Event <- Event{Write, path, info}:
			}
		}
		if oldInfo.Mode() != info.Mode() {
			select {
			case w.Event <- Event{Chmod, path, info}:
			}
		}
	}
}

// for Linux OS check file same
func sameFile(fi1, fi2 os.FileInfo) bool {
	return os.SameFile(fi1, fi2)
}
