package roundtimer

import (
	"fmt"
	"sync"
	"time"
)

type Handle func(rt *RoundTimer)

type RoundTimer struct {
	id               int64
	parameters       interface{}
	interval         time.Duration
	isRound          bool
	running          bool
	timerHandle      Handle
	timerAfterHandle Handle
	timer            *time.Timer

	sync.Mutex
}

func NewRoundTimer() *RoundTimer {
	return &RoundTimer{}
}

// 初次重置计时器回调
func (rt *RoundTimer) initTimer() {
	rt.timer = time.AfterFunc(rt.interval, rt.resetHandle)
	rt.running = true
}

// 重置回调与续期计时器
func (rt *RoundTimer) resetHandle() {
	if !rt.running {
		return
	}
	// 执行定时操作
	rt.timerHandle(rt)
	if !rt.running {
		return
	}
	// 进行续期
	rt.timer = time.AfterFunc(rt.interval, rt.resetHandle)
	if rt.timerAfterHandle != nil {
		// 完成续期的操作
		rt.timerAfterHandle(rt)
	}
}

// 首次启动定时器
func (rt *RoundTimer) Start() error {
	rt.Lock()
	defer rt.Unlock()
	if rt.running {
		return fmt.Errorf("id: %d timer already running", rt.id)
	}
	rt.initTimer()
	return nil
}

func (rt *RoundTimer) StartWithHandle(startHandle Handle) error {
	rt.Lock()
	if rt.running {
		return fmt.Errorf("id: %d timer already running", rt.id)
	}
	rt.initTimer()
	rt.Unlock()
	if startHandle != nil {
		startHandle(rt)
	}
	return nil
}

// 重置定时间隔
func (rt *RoundTimer) ResetInterval(interval time.Duration) error {
	rt.Mutex.Lock()
	defer rt.Mutex.Unlock()

	if !rt.running {
		return fmt.Errorf("id: %d timer already stop", rt.id)
	}

	rt.interval = interval
	rt.timer.Reset(interval)
	return nil
}

func (rt *RoundTimer) Stop() {
	rt.Mutex.Lock()
	defer rt.Mutex.Unlock()

	rt.running = false
	rt.timer.Stop()
}

func (rt *RoundTimer) StopWithHandle(handle Handle) {
	rt.Mutex.Lock()
	rt.running = false
	rt.timer.Stop()
	rt.Mutex.Unlock()
	if handle != nil {
		handle(rt)
	}
}

// 重置对象参数
func (rt *RoundTimer) Reset() {
	rt.Mutex.Lock()
	defer rt.Mutex.Unlock()
	if rt.timer != nil {
		rt.timer.Stop()
	}
	rt.id = 0
	rt.running = false
	rt.interval = 0 * time.Second
	rt.timerHandle = nil
	rt.timerAfterHandle = nil
	rt.parameters = nil
}

func (rt *RoundTimer) SetInterval(interval time.Duration) *RoundTimer {
	if interval < 0 {
		rt.interval = 1 * time.Second
	} else {
		rt.interval = interval
	}
	return rt
}

func (rt *RoundTimer) SetId(id int64) *RoundTimer {
	rt.id = id
	return rt
}

func (rt *RoundTimer) SetTimerHandle(timerHandle func(rt *RoundTimer)) *RoundTimer {
	rt.timerHandle = timerHandle
	return rt
}

func (rt *RoundTimer) SetTimerAfterHandle(timerAfterHandle func(rt *RoundTimer)) *RoundTimer {
	rt.timerAfterHandle = timerAfterHandle
	return rt
}

func (rt *RoundTimer) SetParameters(para interface{}) *RoundTimer {
	rt.parameters = para
	return rt
}

func (rt *RoundTimer) GetPara() interface{} {
	return rt.parameters
}
