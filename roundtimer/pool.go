package roundtimer

import (
	"sync"
	"sync/atomic"
)

// 池化技术: 池化频繁创建的对象
// 复用原有对象, 减少创建和销毁的动作

const DefaultId = 1

var Pool *roundTimerPool

type roundTimerPool struct {
	id int64
	pool *sync.Pool
	mu sync.Mutex
}

// 对象池 避免频繁创建与回收对象
func NewRoundTimerPool() {
	Pool = &roundTimerPool{
		pool: &sync.Pool{
			New: func() interface{} {
				// 创建原始对象
				return NewRoundTimer()
			},
		},
	}
}

// 初始化具备取号
func NewRoundTimerPoolWithId(initId int64) {
	Pool = &roundTimerPool{
		id: initId,
		pool: &sync.Pool{
			New: func() interface{} {
				// 创建原始对象
				return NewRoundTimer()
			},
		},
	}
}

func (r *roundTimerPool) Get() *RoundTimer {
	return r.pool.Get().(*RoundTimer)
}

func (r *roundTimerPool) GetWithId() *RoundTimer {
	r.mu.Lock()
	defer r.mu.Unlock()
	rt := r.pool.Get().(*RoundTimer)
	rt.SetId(r.id)
	atomic.AddInt64(&r.id, 1)
	return rt
}

func (r *roundTimerPool) Put(rt *RoundTimer) {
	rt.Reset()
	r.pool.Put(rt)
}
