package flow_rate

import (
	"math/rand"
	"testing"
	"time"
)

type FlowRateItem struct {
	BookId int64
}

func TestNewFlowRate(t *testing.T) {
	flowRateConfigs := []FlowRateConfig{
		{
			Item:   FlowRateItem{BookId: 62},
			Weight: 60,
		},
		{
			Item:   FlowRateItem{BookId: 70},
			Weight: 10,
		},
		{
			Item:   FlowRateItem{BookId: 71},
			Weight: 5,
		},
		{
			Item:   FlowRateItem{BookId: 72},
			Weight: 25,
		},
	}
	flowRate := NewFlowRate(flowRateConfigs...)
	t.Log(flowRate.printFlowRateWeight())
	testTotal := 1000
	percent := make(map[int64]int)
	for i := 0; i < testTotal; i++ {
		res := flowRate.GetFlowRateResult()
		if item, ok := res.(FlowRateItem); ok {
			percent[item.BookId]++
		}
	}
	for id, cnt := range percent {
		t.Logf("rand-weight: id: %d, percent: %f\n", id, (float64(cnt)/float64(testTotal))*100)
	}

	time.Sleep(time.Second)
	rand.Seed(time.Now().UnixNano())
	flowRate.SetFixedRate()
	percent = make(map[int64]int)
	for i := 0; i < testTotal; i++ {
		res := flowRate.GetFlowRateResult()
		if item, ok := res.(FlowRateItem); ok {
			percent[item.BookId]++
		}
	}
	for id, cnt := range percent {
		t.Logf("fixed-weight: id: %d, percent: %f\n", id, (float64(cnt)/float64(testTotal))*100)
	}
}
