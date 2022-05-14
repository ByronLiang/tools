package fs

import (
	"fmt"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

const path = "data"

var scanTargets = []string{"r1.txt"}

func TestNewFs(t *testing.T) {
	var wg sync.WaitGroup
	fsGroup := make([]*Fs, 0, len(scanTargets))
	for _, target := range scanTargets {
		wg.Add(1)
		fs := NewFs().SetExitHook(InitExitHook(target)).SetScanHandle(InitScanHandle(target))
		fsGroup = append(fsGroup, fs)
		go func(fsObj *Fs, n string) {
			name := filepath.Join(path, n)
			err := fsObj.WithFile(name)
			if err != nil {
				t.Error(err, target)
			}
			wg.Done()
		}(fs, target)
	}
	time.AfterFunc(3*time.Second, func() {
		for _, fs := range fsGroup {
			fs.Exist()
		}
	})
	wg.Wait()
}

func TestFs_SetPosition(t *testing.T) {
	var wg sync.WaitGroup
	fsGroup := make([]*Fs, 0, len(scanTargets))
	for _, target := range scanTargets {
		wg.Add(1)
		// In common case, position should start from 0
		// blew this code the offset 17 (the byte size offset) read from file
		fs := NewFs().
			SetPosition(17).
			SetLine(1).
			SetExitHook(InitExitHook(target)).
			SetScanHandle(InitScanHandle(target))
		fsGroup = append(fsGroup, fs)
		go func(fsObj *Fs, n string) {
			name := filepath.Join(path, n)
			err := fsObj.WithFile(name)
			if err != nil {
				t.Error(err, target)
			}
			wg.Done()
		}(fs, target)
	}
	wg.Wait()
}

func InitExitHook(name string) ExitHook {
	hook := func(position, line int64) {
		fmt.Printf("exit file scaner name: [%s] pos is %d, line: %d \n", name, position, line)
		return
	}
	return hook
}

func InitScanHandle(name string) ScanHandle {
	return func(content string, position, line int64) {
		fmt.Printf("name: [%s] content: [%s] pos is %d, line: %d \n", name, content, position, line)
		time.Sleep(800 * time.Millisecond)
		return
	}
}
