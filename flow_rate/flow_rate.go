package flow_rate

import (
	"fmt"
	"math/rand"
	"time"
)

// 流量切割权重配置
type FlowRateConfig struct {
	Item   interface{}
	Weight int // 权重
}

type FlowRateWeight struct {
	flowRateItemIndex int // get data from FlowRateItems
	l, r              int // l: 左边界 r: 右边界 eg: [l,r)
}

type RateStrategy interface {
	GetRate() int
}

type FlowRate struct {
	RateStrategy
	FlowRateItems   []interface{}
	FlowRateWeights []FlowRateWeight
	sumWeight       int
}

func NewFlowRate(flowRateConfigs ...FlowRateConfig) *FlowRate {
	rand.Seed(time.Now().UnixNano())
	flowRate := &FlowRate{
		FlowRateItems:   make([]interface{}, 0, len(flowRateConfigs)),
		FlowRateWeights: make([]FlowRateWeight, 0, len(flowRateConfigs)),
	}
	var rSum, lSum int
	for i, flowRateCfg := range flowRateConfigs {
		rSum += flowRateCfg.Weight
		flowRateWeight := FlowRateWeight{
			flowRateItemIndex: i,
			r:                 rSum,
		}
		if i != 0 {
			lSum += flowRateConfigs[i-1].Weight
			flowRateWeight.l = lSum
		}
		flowRate.FlowRateItems = append(flowRate.FlowRateItems, flowRateCfg.Item)
		flowRate.FlowRateWeights = append(flowRate.FlowRateWeights, flowRateWeight)
	}
	// 填充总权重值
	flowRate.sumWeight = rSum
	flowRate.SetRandRate()
	return flowRate
}

func (f *FlowRate) SetFixedRate() *FlowRate {
	f.RateStrategy = NewFlowFixedRate(f.sumWeight)
	return f
}

func (f *FlowRate) SetRandRate() *FlowRate {
	f.RateStrategy = NewFlowRandRate(f.sumWeight)
	return f
}

// 获取流量切割对象
func (f *FlowRate) GetFlowRateResult() interface{} {
	randWeight := f.RateStrategy.GetRate()
	for _, flowRateWeight := range f.FlowRateWeights {
		if flowRateWeight.l <= randWeight && randWeight < flowRateWeight.r {
			return f.FlowRateItems[flowRateWeight.flowRateItemIndex]
		}
	}
	return nil
}

func (f *FlowRate) printFlowRateWeight() string {
	var str string
	for _, flowRateWeight := range f.FlowRateWeights {
		str += fmt.Sprintf("index: %d, range: l %d, r %d\n", flowRateWeight.flowRateItemIndex, flowRateWeight.l, flowRateWeight.r)
	}
	return str
}
