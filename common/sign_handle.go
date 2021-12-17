package common

import (
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func SignCatch(quitHandle func(content []byte), intHandle func()) {
	signs := make(chan os.Signal, 1)
	signal.Notify(signs, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGKILL)
	switch <-signs {
	// quit异常退出信号捕捉 打印stack 终端里：Ctrl 与 \ 发出此信号
	case syscall.SIGQUIT:
		buf := make([]byte, 1<<20)
		slacken := runtime.Stack(buf, true)
		if quitHandle == nil {
			defaultQuitHandle(buf[:slacken])
		} else {
			quitHandle(buf[:slacken])
		}
	case syscall.SIGINT:
		if intHandle == nil {
			defaultIntHandle()
		} else {
			intHandle()
		}
	case syscall.SIGKILL:
		if intHandle == nil {
			defaultIntHandle()
		} else {
			intHandle()
		}
	}
}

func defaultQuitHandle(content []byte) {
	log.Printf("=== received SIGQUIT ===\n*** goroutine dump...\n%s\n*** end\n", content)
}

func defaultIntHandle() {
	log.Println("== SIGINT/SIGKILL ==")
}
