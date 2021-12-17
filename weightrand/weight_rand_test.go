package weightrand

import (
	"testing"
)

type RoomItem struct {
	Weight int64
	RoomId int64
}

func initWeightItems() []*weightRandItem {
	return nil
}

func TestNewWeightRand(t *testing.T) {
	obj := NewWeightRand(
		NewWeightRandItem(1, 1),
		NewWeightRandItem(2, 2),
		NewWeightRandItem(1, 3)).InitTotal().CalMax()
	obj.InitStrategy(&CommonWeightRand{Wr: obj})
	t.Logf("total: %v; max: %v", obj.GetTotal(), obj.GetMax())
	res, err := obj.Pick()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("res: %v", res)
}

func TestWeightRand_Pick(t *testing.T) {
	obj := NewWeightRand(
		NewWeightRandItem(1, 1),
		NewWeightRandItem(2, 2),
		NewWeightRandItem(1, 3)).CalMax()
	obj.InitStrategy(&SubWeightRand{Wr: obj})
	res, err := obj.Pick()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("res: %v", res)
}

func TestWeightRand_MulPick(t *testing.T) {
	obj := NewWeightRand(
		NewWeightRandItem(1, 1),
		NewWeightRandItem(2, 2),
		NewWeightRandItem(1, 3)).CalMax()
	obj.InitStrategy(&SubWeightRand{Wr: obj})
	res, err := obj.MulPick(5)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("res: %v", res)
}

func TestWeightRand_BinaryPick(t *testing.T) {
	obj := NewWeightRand(
		NewWeightRandItem(1, 1),
		NewWeightRandItem(2, 2),
		NewWeightRandItem(1, 3)).InitTotal().CalMax()
	obj.InitStrategy(&BinaryWeightRand{Wr: obj})
	res, err := obj.MulPick(5)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("res: %v", res)
}
