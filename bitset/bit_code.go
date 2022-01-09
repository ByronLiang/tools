package bitset

import (
	"errors"
	"fmt"
	"sync"
)

// 最大随机码可存储数据位数63
// 若存储64位, 则将uint64 转化为 int64 会丢失最高位的数据 int64 需要有一个数位是符号位
const MaxBitLen = 63

var (
	OverBitLenErr     = errors.New("over max bit len 63")
	GroupKeyNoFindErr = errors.New("key no exist")
)

type bitCodeElement struct {
	Key        string `json:"key"`
	BitLen     int    `json:"bit_len"`
	bitPos     uint64 `json:"-"` // 当前元素字节位置
	restBitPos uint64 `json:"-"` // 剩余元素字节位数
}

type BitCodeConfig struct {
	Elements []*bitCodeElement
	ByKey    map[string]int
}

type BitCodeGroup struct {
	mu     sync.RWMutex
	group  map[string]*BitCodeConfig
	bitLen int
}

func NewBitCodeGroup() *BitCodeGroup {
	return &BitCodeGroup{
		mu:     sync.RWMutex{},
		group:  make(map[string]*BitCodeConfig),
		bitLen: MaxBitLen,
	}
}

// 初始化配置并加入组里
func (bcg *BitCodeGroup) Add(groupKey string, elements ...*bitCodeElement) error {
	bcg.mu.Lock()
	defer bcg.mu.Unlock()
	byKeySet := make(map[string]int, len(elements))
	paddingBitLen := 0
	for i, element := range elements {
		// 当前元素位数
		element.bitPos = uint64(1<<(bcg.bitLen-paddingBitLen) - 1)
		// 剩余元素的位数
		element.restBitPos = uint64(bcg.bitLen - paddingBitLen - element.BitLen)
		// 当前元素已提取, 进行偏移处理
		paddingBitLen += element.BitLen
		// 标记 key 与 下标
		byKeySet[element.Key] = i
	}
	// 校验是否超出位数
	if paddingBitLen > bcg.bitLen {
		return OverBitLenErr
	}
	bcg.group[groupKey] = &BitCodeConfig{Elements: elements, ByKey: byKeySet}
	return nil
}

func (bcg *BitCodeGroup) Delete(groupKey string) error {
	bcg.mu.Lock()
	defer bcg.mu.Unlock()
	_, ok := bcg.group[groupKey]
	if !ok {
		return GroupKeyNoFindErr
	}
	delete(bcg.group, groupKey)
	return nil
}

// 按照元素生成随机码
func (bcg *BitCodeGroup) GenerateCode(groupKey string, elements map[string]uint64) (uint64, error) {
	bcg.mu.RLock()
	defer bcg.mu.RUnlock()
	code := uint64(0)
	bitCodeCfg, ok := bcg.group[groupKey]
	if !ok {
		return code, GroupKeyNoFindErr
	}
	paddingBitLen := 0
	for _, element := range bitCodeCfg.Elements {
		value, ok := elements[element.Key]
		if !ok {
			paddingBitLen += element.BitLen
			continue
		}
		// 校验是否超出 bitLen
		if uint64(1<<element.BitLen-1) < value {
			return code, fmt.Errorf("key: %s element value over bit len: %d", element.Key, element.BitLen)
		}
		paddingBitLen += element.BitLen
		code |= uint64(value << (bcg.bitLen - paddingBitLen))
	}
	return code, nil
}

// 从随机号解析指定参数
func (bcg *BitCodeGroup) Parse(groupKey string, code uint64, elements ...string) (map[string]uint64, error) {
	bcg.mu.RLock()
	defer bcg.mu.RUnlock()
	bitCodeCfg, ok := bcg.group[groupKey]
	if !ok {
		return nil, GroupKeyNoFindErr
	}
	res := make(map[string]uint64, len(elements))
	for _, key := range elements {
		i, ok := bitCodeCfg.ByKey[key]
		if !ok {
			continue
		}
		res[key] = code & bitCodeCfg.Elements[i].bitPos >> bitCodeCfg.Elements[i].restBitPos
	}
	return res, nil
}

// 解析随机号全部参数
func (bcg *BitCodeGroup) ParseAll(groupKey string, code uint64) (map[string]uint64, error) {
	bcg.mu.RLock()
	defer bcg.mu.RUnlock()
	bitCodeCfg, ok := bcg.group[groupKey]
	if !ok {
		return nil, GroupKeyNoFindErr
	}
	res := make(map[string]uint64, len(bitCodeCfg.Elements))
	for _, element := range bitCodeCfg.Elements {
		res[element.Key] = code & element.bitPos >> element.restBitPos
	}
	return res, nil
}

func NewBitCodeElement(key string, bitLen int) *bitCodeElement {
	return &bitCodeElement{
		Key:    key,
		BitLen: bitLen,
	}
}
