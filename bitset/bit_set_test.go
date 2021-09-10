package bitset

import "testing"

const DefaultBit BitUnit = iota

// BitUnit是八位，只能支持7个bit位配置, 其中一位是符号位
const (
	ConfA BitUnit = 1 << iota
	ConfB
	ConfC
)

type ConfItem struct {
	Key   string `json:"key"`
	OnOff bool   `json:"onOff"`
}

func TestAddBit(t *testing.T) {
	nowBit := AddBit(DefaultBit, ConfB)
	nowBit = AddBit(nowBit, ConfC)
	t.Log(BitString(nowBit), nowBit)

	nowBit = DelBit(nowBit, ConfB)
	t.Log(BitString(nowBit), nowBit)

	configMap := make(map[BitUnit]*ConfItem)
	configMap[ConfA] = &ConfItem{Key: "confA"}
	configMap[ConfB] = &ConfItem{Key: "confB"}
	configMap[ConfC] = &ConfItem{Key: "confC"}

	for conf, item := range configMap {
		item.OnOff = Exist(nowBit, conf)
	}

	confList := []BitUnit{ConfA, ConfB, ConfC}
	for _, conf := range confList {
		t.Log("exist: ", Exist(nowBit, conf))
	}
}

func TestBitString(t *testing.T) {
	// 因溢出而报错
	const OverConf = 1 << 6
	nowBit := AddBit(DefaultBit, OverConf)
	t.Log(BitString(nowBit), nowBit)
}
