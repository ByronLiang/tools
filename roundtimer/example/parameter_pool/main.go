package main

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/ByronLiang/tools/common"
	"github.com/ByronLiang/tools/roundtimer"
)

func main() {
	roundtimer.NewRoundTimerPool()
	rt := roundtimer.Pool.Get()
	rt.SetInterval(1 * time.Second).
		SetId(1).
		SetParameters(&CountPara{
			Round: 0,
			Total: 10,
		}).
		SetTimerHandle(paraCallback)
	err := rt.Start()
	if err != nil {
		fmt.Println(err)
		return
	}
	common.SignCatch(nil, func() {
		roundtimer.Pool.Put(rt)
	})
}

type CountPara struct {
	Round int64
	Total int64
}

// 复杂定时器任务: 动态调整定时间隔
func paraCallback(rt *roundtimer.RoundTimer) {
	rt.Lock()
	defer rt.Unlock()
	if para, ok := rt.GetPara().(*CountPara); ok {
		if para.Total-para.Round == 0 {
			rt.SetInterval(2 * time.Second)
			rt.SetParameters(&CountPara{
				Round: 0,
				Total: 5,
			})
		} else {
			if para.Round == 3 {
				fmt.Println("round 3 counter")
			} else {
				fmt.Printf("timer count: %d \n", para.Total-para.Round)
			}
			atomic.AddInt64(&para.Round, 1)
			rt.SetParameters(para)
		}
		if para.Total == 5 && para.Round == 5 {
			rt.StopWithHandle(func(_ *roundtimer.RoundTimer) {
				fmt.Println("end the timer")
			})
		}
	}
}
