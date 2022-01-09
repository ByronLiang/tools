package bitset

import (
	"fmt"
	"testing"
	"time"
)

func testBuildBitCodeElement() []*bitCodeElement {
	ele1 := NewBitCodeElement("shopNo", 4)
	ele2 := NewBitCodeElement("countryNo", 4)
	ele3 := NewBitCodeElement("timestamp", 32)
	return []*bitCodeElement{ele1, ele2, ele3}
}

func TestBitCodeGroup_Add(t *testing.T) {
	bcg := NewBitCodeGroup()
	ele1 := NewBitCodeElement("shopNo", 4)
	ele2 := NewBitCodeElement("countryNo", 34)
	ele3 := NewBitCodeElement("timestamp", 32)
	err := bcg.Add("bit-err", ele1, ele2, ele3)
	if err != nil {
		t.Error(err)
	}
}

func TestBitCodeGroup_GenerateCodeErr(t *testing.T) {
	bcg := NewBitCodeGroup()
	elem := testBuildBitCodeElement()
	err := bcg.Add("bit-element-over-err", elem...)
	if err != nil {
		t.Fatal(err)
	}
	data := map[string]uint64{
		"shopNo":    uint64(2),
		"countryNo": uint64(20),
		"timestamp": uint64(time.Now().Unix()),
	}
	code, err := bcg.GenerateCode("bit-element-over-err", data)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(code)
}

func TestNewBitCodeGroup(t *testing.T) {
	bcg := NewBitCodeGroup()
	elem := testBuildBitCodeElement()
	err := bcg.Add("orderNo", elem...)
	if err != nil {
		t.Fatal(err)
	}
	codeList := make(map[uint64]map[string]uint64)
	for i := 0; i <= 3; i++ {
		data := map[string]uint64{
			"shopNo":    uint64(i + 1),
			"countryNo": uint64(i + 2),
			"timestamp": uint64(time.Now().Unix()),
		}
		time.Sleep(time.Duration(i) * time.Second)
		code, err := bcg.GenerateCode("orderNo", data)
		if err != nil {
			t.Error(err)
			return
		}
		codeList[code] = data
	}
	fmt.Println(codeList)
	time.Sleep(1 * time.Second)
	for code, _ := range codeList {
		res, err := bcg.ParseAll("orderNo", code)
		if err != nil {
			t.Error(err)
			break
		}
		fmt.Println(res)
		parseMap, err := bcg.Parse("orderNo", code, "countryNo", "shopNo")
		if err != nil {
			t.Error(err)
		} else {
			fmt.Println("parseMap", parseMap)
		}
	}
}

func TestBitCodeGroup_Delete(t *testing.T) {
	bcg := NewBitCodeGroup()
	elem := testBuildBitCodeElement()
	err := bcg.Add("orderNo", elem...)
	if err != nil {
		t.Fatal(err)
	}
	elem2 := testBuildBitCodeElement()
	err = bcg.Add("paperNo", elem2...)
	if err != nil {
		t.Fatal(err)
	}
	err = bcg.Delete("orderNo")
	if err != nil {
		t.Fatal(err)
	}
	data := map[string]uint64{
		"shopNo":    uint64(1),
		"countryNo": uint64(2),
		"timestamp": uint64(time.Now().Unix()),
	}
	code, err := bcg.GenerateCode("paperNo", data)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(code)
	code, err = bcg.GenerateCode("orderNo", data)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(code)
}
