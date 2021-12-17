package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/ByronLiang/tools/daemon"
)

// 需要将程序进行包编译，执行./包名 -d 程序将在后台运行
func main() {
	daemon.Daemon(handle)
}

func handle(c chan struct{}) {
	file, err := os.OpenFile("dae_test.txt", os.O_CREATE|os.O_RDWR, 0664)
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()
	for {
		select {
		case <-c:
			return
		default:
			_, err := file.Write([]byte(strconv.Itoa((int)(time.Now().Unix())) + "\n"))
			if err != nil {
				return
			}
			time.Sleep(time.Second * 1)
		}
	}
}
