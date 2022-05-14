package fs

import (
	"bufio"
	"io"
	"os"
)

type ScanHandle func(content string, position, line int64)

type ExitHook func(position, line int64)

type Fs struct {
	position   int64 // position for seek start file byte
	line       int64
	exit       chan struct{}
	scanHandle ScanHandle
	exitHook   ExitHook
}

func NewFs() *Fs {
	return &Fs{
		exit: make(chan struct{}, 1),
	}
}

func (fs *Fs) SetPosition(position int64) *Fs {
	fs.position = position
	return fs
}

func (fs *Fs) SetLine(line int64) *Fs {
	fs.line = line
	return fs
}

func (fs *Fs) SetScanHandle(handle ScanHandle) *Fs {
	fs.scanHandle = handle
	return fs
}

func (fs *Fs) SetExitHook(hook ExitHook) *Fs {
	fs.exitHook = hook
	return fs
}

func (fs *Fs) WithFile(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	err = fs.WithScanner(f)
	return err
}

// save position with scanLines
func (fs *Fs) WithScanner(input io.ReadSeeker) error {
	// 定位起始点
	if _, err := input.Seek(fs.position, io.SeekStart); err != nil {
		return err
	}
	scanner := bufio.NewScanner(input)
	scanner.Split(fs.splitHandle)
	for scanner.Scan() {
		select {
		case <-fs.exit:
			if fs.exitHook != nil {
				fs.exitHook(fs.position, fs.line)
			}
			return scanner.Err()
		default:
			fs.line++
			if fs.scanHandle != nil {
				fs.scanHandle(scanner.Text(), fs.position, fs.line)
			}
		}
	}
	return scanner.Err()
}

func (fs *Fs) GetPosition() int64 {
	return fs.position
}

func (fs *Fs) GetLine() int64 {
	return fs.line
}

func (fs *Fs) Exist() {
	fs.exit <- struct{}{}
}

func (fs *Fs) splitHandle(data []byte, atEOF bool) (advance int, token []byte, err error) {
	advance, token, err = bufio.ScanLines(data, atEOF)
	fs.position += int64(advance)
	return
}
