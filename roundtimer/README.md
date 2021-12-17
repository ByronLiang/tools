# 自定义定时组件

## 自定义回调解析

1. 定时规则持久化处理：由于无法对`func`进行序列化；`json: unsupported type: func()` 需要对回调指令, 进行二次解析

2. 若无法进行持久化, 只能以进程缓存进行处理数据

### 基本回调需求解析

1. 到达指定秒数, 相关内容广播，接口调用等

## 潜在线程不安全设计

```go
package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	start := time.Now()
	var t *time.Timer
	t = time.AfterFunc(randomDuration(), func() {
		fmt.Println(time.Now().Sub(start))
		t.Reset(randomDuration())
	})
	time.Sleep(5 * time.Second)
}

func randomDuration() time.Duration {
	return time.Duration(rand.Int63n(1e9))
}
```

因`main`协程与`time.AfterFunc`的协程函数因调用的先后顺序，会引发全局变量`t *time.Timer`无法被赋值，`t.Reset(randomDuration())`会因空指针诱发异常

### 检测

执行单元测试时: `go test -race app(包名)`

执行运行时：`go run -race main.go`

### 避免协程竞争

不使用`Reset(d Duration)`进行重置, 重新初始化定时与定时回调：`AfterFunc(d Duration, f func())`
