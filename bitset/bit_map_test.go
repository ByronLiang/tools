package bitset

import (
	"fmt"
	"testing"
)

func TestNewBitMap(t *testing.T) {
	bm := NewBitMap()
	err := bm.Add(5)
	if err != nil {
		t.Fatal(err)
		return
	}
	err = bm.Add(13)
	if err != nil {
		t.Log(err)
	}
	if bm.IsExist(10) {
		t.Log("num 10 exist")
	}
	if bm.IsExist(12) {
		t.Log("num 12 exist")
	}
	err = bm.Remove(90)
	if err != nil {
		t.Logf("remove 90 err: %s", err)
	}
	bm.Remove(5)
	// add
	_ = bm.Add(3)
	_ = bm.Add(63)
	_ = bm.Add(64)
	err = bm.Add(13)
	if err != nil {
		t.Logf("add num 13 err: %s", err)
	}
	fmt.Println(bm.BitMapString())
	fmt.Println(bm.Values())
	fmt.Println(bm.Len())
}
