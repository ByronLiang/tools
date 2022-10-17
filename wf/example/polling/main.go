package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ByronLiang/tools/wf"
)

func main() {
	polling := wf.NewPolling(15 * time.Second)
	// temp file
	tempFileName, err := wf.BuildTempFile("polling_*.txt")
	if err != nil {
		log.Println("tempFile", err)
		return
	}
	log.Println("new-add: ", tempFileName)
	err = polling.Add(tempFileName)
	if err != nil {
		log.Println("polling-add", err)
		return
	}
	// polling dir
	err = polling.Add("/tmp/polling_dir")
	if err != nil {
		log.Println("polling-add", err)
		return
	}
	go func() {
		for {
			select {
			case event := <-polling.Event:
				if event.Op&(wf.Write|wf.Create) != 0 {
					log.Println("event:", event, event.Op)
				}
			case pollingErr := <-polling.Error:
				log.Println("pollingErr", pollingErr)
			}
		}
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGKILL}...)
	select {
	case <-c:
		log.Println("end-watch")
		polling.Close()
	}
	time.Sleep(1 * time.Second)
}

func mockAddFile(duration time.Duration, polling wf.Watcher) {
	// 模拟新增文件
	time.AfterFunc(duration, func() {
		tempFileName, err := wf.BuildTempFile("polling_*.txt")
		if err != nil {
			log.Println("tempFile", err)
			return
		}
		polling.Add(tempFileName)
	})
}
