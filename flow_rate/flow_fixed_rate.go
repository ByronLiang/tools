package flow_rate

import (
	"math/rand"
	"sync/atomic"
)

type FlowFixedRate struct {
	base     int
	op       int
	atomicOp uint64 // 原子性计数器
	source   []int
}

func NewFlowFixedRate(sumWeight int) *FlowFixedRate {
	flowFixRate := &FlowFixedRate{base: sumWeight}
	source := make([]int, flowFixRate.base, flowFixRate.base)
	for i := 0; i < flowFixRate.base; i++ {
		source[i] = i
	}
	// 随机排序
	rand.Shuffle(flowFixRate.base, func(i, j int) {
		source[i], source[j] = source[j], source[i]
	})
	flowFixRate.source = source

	return flowFixRate
}

func (f *FlowFixedRate) GetRate() int {
	// 原子性
	return f.source[int(atomic.AddUint64(&f.atomicOp, 1))%f.base]
	// 非原子性
	//f.op++
	//f.op = f.op % f.base
	//return f.source[f.op]
}
