package weightrand

import (
	cr "crypto/rand"
	"errors"
	"math/big"
	"math/rand"
	"sort"
	"time"
)

/**
基于权重随机包
Weighted random generation
*/
type WeightRand struct {
	members  []*weightRandItem
	total    []int64
	max      int64
	strategy RandStrategy
}

type weightRandItem struct {
	Weight int64
	Item   interface{}
}

func NewWeightRandItem(weight int64, item interface{}) *weightRandItem {
	return &weightRandItem{
		Weight: weight,
		Item:   item,
	}
}

func NewWeightRand(items ...*weightRandItem) *WeightRand {
	return &WeightRand{
		members: items,
		total:   nil,
	}
}

func (wr *WeightRand) CalMax() *WeightRand {
	total := int64(0)
	for _, item := range wr.members {
		total += item.Weight
	}
	wr.max = total
	return wr
}

func (wr *WeightRand) SortMembers() *WeightRand {
	// 成员按照权重排序
	sort.Slice(wr.members, func(i, j int) bool {
		return wr.members[i].Weight < wr.members[j].Weight
	})
	return wr
}

func (wr *WeightRand) InitTotal() *WeightRand {
	total := make([]int64, 0, len(wr.members))
	plusTotal := int64(0)
	for _, item := range wr.members {
		plusTotal += item.Weight
		total = append(total, plusTotal)
	}
	wr.total = total
	return wr
}

func (wr *WeightRand) GetTotal() []int64 {
	return wr.total
}

func (wr *WeightRand) GetMax() int64 {
	return wr.max
}

func (wr *WeightRand) InitStrategy(strategy RandStrategy) *WeightRand {
	wr.strategy = strategy

	return wr
}

func (wr *WeightRand) Pick() (interface{}, error) {
	seed := InitRand(wr.max)
	result := wr.strategy.PickMember(seed)
	if result == nil {
		return nil, errors.New("un match")
	}
	return result, nil
}

func (wr *WeightRand) MulPick(num int) ([]interface{}, error) {
	if num <= 0 {
		return nil, errors.New("illegal num")
	}
	data := make([]interface{}, 0, num)
	for num > 0 {
		seed := InitRand(wr.max)
		result := wr.strategy.PickMember(seed)
		if result != nil {
			data = append(data, result)
		}
		num--
	}
	return data, nil
}

func genRand(max int64) int64 {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Int63n(max)
}

func InitRand(max int64) int64 {
	if res, err := cr.Int(cr.Reader, big.NewInt(max)); err == nil {
		return res.Int64()
	}
	return 0
}
