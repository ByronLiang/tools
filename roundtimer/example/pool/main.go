package main

import (
	"fmt"
	"time"

	"github.com/ByronLiang/tools/common"
	"github.com/ByronLiang/tools/roundtimer"
)

func main() {
	roundtimer.NewRoundTimerPool()
	rt := roundtimer.Pool.Get()
	rt.SetInterval(2 * time.Second).
		SetId(1).
		SetTimerHandle(roundTimerCallback).
		SetTimerAfterHandle(roundTimerResetAfter)
	err := rt.Start()
	if err != nil {
		fmt.Println(err)
		return
	}
	common.SignCatch(nil, func() {
		roundtimer.Pool.Put(rt)
	})
}

func roundTimerCallback(rt *roundtimer.RoundTimer) {
	fmt.Println("timer count")
}

func roundTimerResetAfter(rt *roundtimer.RoundTimer) {
	fmt.Println("after reset timer count")
}
