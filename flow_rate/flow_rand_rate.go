package flow_rate

import (
	"math/rand"
	"time"
)

type FlowRandRate struct {
	sumWeight int
}

func NewFlowRandRate(sumWeight int) *FlowRandRate {
	rand.Seed(time.Now().UnixNano())
	return &FlowRandRate{sumWeight: sumWeight}
}

func (f *FlowRandRate) GetRate() int {
	return rand.Intn(f.sumWeight)
}
