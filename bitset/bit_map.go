package bitset

import (
	"errors"
	"strconv"
	"sync"
)

var (
	ErrNumHadExist = errors.New("number had exist")
	ErrNumNoExist  = errors.New("number no exist")
)

type BitMap struct {
	mu sync.RWMutex
	//map key: num / 64; map value: num % 64
	modValues map[int64]uint64
	length    int64 // bit map element size
	modUnit   int64 // mod unit 取模值
}

func NewBitMap() *BitMap {
	return &BitMap{
		mu:        sync.RWMutex{},
		modValues: make(map[int64]uint64),
		length:    0,
		modUnit:   64,
	}
}

func (bm *BitMap) IsExist(num int64) bool {
	index := num / bm.modUnit
	modBit := uint(num % bm.modUnit)
	bm.mu.RLock()
	defer bm.mu.RUnlock()
	if modValue, ok := bm.modValues[index]; ok {
		return modValue&(1<<modBit) != 0
	}
	return false
}

func (bm *BitMap) Add(num int64) error {
	index := num / bm.modUnit
	modBit := uint(num % bm.modUnit)
	bm.mu.Lock()
	defer bm.mu.Unlock()
	if modValue, ok := bm.modValues[index]; ok {
		if modValue&(1<<modBit) != 0 {
			return ErrNumHadExist
		}
	}
	bm.modValues[index] |= uint64(1 << modBit)
	bm.length++
	return nil
}

func (bm *BitMap) Remove(num int64) error {
	index := num / bm.modUnit
	modBit := uint(num % bm.modUnit)
	bm.mu.Lock()
	defer bm.mu.Unlock()
	if modValue, ok := bm.modValues[index]; ok {
		if modValue&(1<<modBit) == 0 {
			return ErrNumNoExist
		}
		bm.modValues[index] = modValue - uint64(1<<modBit)
		bm.length--
		return nil
	}
	return ErrNumNoExist
}

func (bm *BitMap) Len() int64 {
	return bm.length
}

// 将位图存放的数值全部转换成十进制
func (bm *BitMap) Values() map[int64]struct{} {
	bm.mu.RLock()
	defer bm.mu.RUnlock()
	res := make(map[int64]struct{})
	for i, modBit := range bm.modValues {
		prefix := i * bm.modUnit
		for j := int64(0); j < bm.modUnit; j++ {
			if modBit&(1<<j) != 0 {
				res[prefix+j] = struct{}{}
			}
		}
	}
	return res
}

// 显示二进制字符串
func (bm *BitMap) BitMapString() map[int64]string {
	bm.mu.RLock()
	defer bm.mu.RUnlock()
	res := make(map[int64]string)
	for index, modValue := range bm.modValues {
		modValueStr := strconv.FormatUint(modValue, 2)
		modValueByte := []byte(modValueStr)
		i, j := 0, len(modValueByte)-1
		for i < j {
			modValueByte[j], modValueByte[i] = modValueByte[i], modValueByte[j]
			i++
			j--
		}
		res[index] = string(modValueByte)
	}
	return res
}
