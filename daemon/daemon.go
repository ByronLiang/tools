package daemon

import (
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

func Daemon(handle func(c chan struct{})) {
	args := os.Args
	daemon := false
	for k, v := range args {
		if v == "-d" {
			daemon = true
			args[k] = ""
		}
	}
	if daemon {
		// fork子进程
		daemonHandle(args...)
		// 父进程结束
		return
	}
	// 相关信号监听, 只能手动kill进程号
	endSign := make(chan struct{})
	signs := make(chan os.Signal, 1)
	signal.Notify(signs, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGKILL)
	go handle(endSign)
	<-signs
	endSign <- struct{}{}
	time.Sleep(1 * time.Second)
}

// 通过系统调用产生子进程
func daemonHandle(args ...string) {
	var arg []string
	if len(args) > 1 {
		arg = args[1:]
	}
	cmd := exec.Command(args[0], arg...)
	cmd.Env = os.Environ()
	cmd.Start()
}
