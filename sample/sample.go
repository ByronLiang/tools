package sample

import (
	"sync/atomic"
	"time"
)

type RefreshCounter struct {
	resetAt int64
	counter int64
}

func (c *RefreshCounter) IsRefresh(t time.Time, tick time.Duration, maxTotal int64) (isRefresh bool) {
	var nowCounter int64
	tn := t.UnixNano()
	resetAfter := atomic.LoadInt64(&c.resetAt)
	// 指定时间范围内, 累加计数器
	if resetAfter > tn {
		nowCounter = atomic.AddInt64(&c.counter, 1)
		// 未超出最大值
		if nowCounter < maxTotal {
			return
		}
	}
	// 针对初始化情况, 获取当前counter 值
	if nowCounter == 0 {
		nowCounter = atomic.LoadInt64(&c.counter)
	}
	// 超出时间范围/超出最大值
	// 重置起始值为1
	atomic.StoreInt64(&c.counter, 1)
	// 重置时间有效期
	newResetAfter := tn + tick.Nanoseconds()
	// CAS失败, 返回起始值
	if !atomic.CompareAndSwapInt64(&c.resetAt, resetAfter, newResetAfter) {
		// 对本次进行累加计算
		isRefresh = false
		return
	}
	// 针对初始化情况 nowCounter = 1 不进行重刷
	if nowCounter > 1 {
		isRefresh = true
	}
	return
}

func (c *RefreshCounter) IncCheckReset(t time.Time, tick time.Duration) int64 {
	tn := t.UnixNano()
	resetAfter := atomic.LoadInt64(&c.resetAt)
	// 指定时间范围内, 累加计数器
	if resetAfter > tn {
		return atomic.AddInt64(&c.counter, 1)
	}
	// 超出时间范围: 取出原累加数值
	nowCounter := atomic.LoadInt64(&c.counter)
	// 重置起始值为1
	atomic.StoreInt64(&c.counter, 1)
	// 重置时间有效期
	newResetAfter := tn + tick.Nanoseconds()
	// CAS失败, 返回起始值
	if !atomic.CompareAndSwapInt64(&c.resetAt, resetAfter, newResetAfter) {
		// 对本次进行累加计算
		return atomic.AddInt64(&c.counter, 1)
	}
	// 针对初始化情况, 初始值为0
	if nowCounter == 0 {
		return 1
	}
	return nowCounter
}
