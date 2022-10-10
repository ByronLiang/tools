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
	watchFile := wf.NewWf()
	if watchFile == nil {
		return
	}
	defer watchFile.Close()
	// 挂载目录
	notifyFile, err := wf.GenWatchPath("/data/fs", "fw.txt")
	if err != nil {
		return
	}
	// 镜像系统路径
	tickerFile, err := wf.GenWatchPath("/tmp/ticker", "check.txt")
	if err != nil {
		return
	}
	autoFileModifyCheck, err := wf.InitAutoFileModify(tickerFile, 60*60)
	if err != nil {
		return
	}
	if watchFile.SetHeartBeat(autoFileModifyCheck) != nil {
		return
	}
	if watchFile.AddFile(notifyFile) != nil {
		return
	}
	watchFile.ToWatch()
	c := make(chan os.Signal, 1)
	signal.Notify(c, []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGKILL}...)
	select {
	case <-c:
		log.Println("total", watchFile.GetCounter())
		log.Println("end-watch")
		time.Sleep(1 * time.Second)
	}
}
