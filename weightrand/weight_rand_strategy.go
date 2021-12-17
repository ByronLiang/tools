package weightrand

type RandStrategy interface {
	PickMember(seed int64) interface{}
}

type CommonWeightRand struct {
	Wr *WeightRand
}

/**
权重累加, 对随机数进行定位并得到下标索引值
*/
func (cw *CommonWeightRand) PickMember(seed int64) interface{} {
	for index, weight := range cw.Wr.total {
		if seed < weight {
			return cw.Wr.members[index].Item
		}
	}
	return nil
}

type SubWeightRand struct {
	Wr *WeightRand
}

/**
对随机值进行减操作, 无需对权重值进行递加处理
*/
func (sw *SubWeightRand) PickMember(seed int64) interface{} {
	for index, member := range sw.Wr.members {
		if seed < member.Weight {
			return sw.Wr.members[index].Item
		}
		seed -= member.Weight
	}
	return nil
}

type BinaryWeightRand struct {
	Wr *WeightRand
}

/**
二分法查询
*/
func (bw *BinaryWeightRand) PickMember(seed int64) interface{} {
	i, j := 0, len(bw.Wr.total)-1
	for i < j {
		mid := (i + (j - i)) >> 1
		if seed >= bw.Wr.total[mid] {
			i = mid + 1
		} else {
			j = mid
		}
	}
	return bw.Wr.members[i].Item
}
