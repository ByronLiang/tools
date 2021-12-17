package common

import (
	"fmt"
	"strconv"
	"sync"
)

/**
https://leetcode-cn.com/problems/excel-sheet-column-title/
*/

type CodeGen struct {
	mu       sync.Mutex
	InitCode int
	// 数字后缀长度转换
	numWeight int
	// 字母前缀
	PrefixSize int
	// 数字后缀字符长度
	NumLength int
}

func NewCodeGen(code, prefixSize, numberLength int) *CodeGen {
	cg := &CodeGen{
		InitCode:   code,
		PrefixSize: prefixSize,
		NumLength:  numberLength,
	}
	cg.initNumWeight()
	return cg
}

func (cg *CodeGen) ReflectCode() string {
	n := cg.InitCode / cg.numWeight
	prefix := cg.genPrefix(n)

	return fmt.Sprintf("%s%s", prefix, strconv.Itoa(cg.InitCode%cg.numWeight))
}

func (cg *CodeGen) genPrefix(code int) string {
	prefixByte := make([]rune, cg.PrefixSize)
	index := cg.PrefixSize - 1
	for code > 0 {
		m := code % 26
		if m == 0 {
			m = 26
			code -= 1
		}
		code = code / 26
		prefixByte[index] = rune(m + 64)
		index--
	}
	return string(prefixByte)
}

func (cg *CodeGen) GenCode() string {
	cg.mu.Lock()
	var codeIndex = cg.InitCode
	cg.InitCode += 1
	cg.mu.Unlock()
	initPrefix := codeIndex / cg.numWeight
	initNum := codeIndex % cg.numWeight
	prefix := cg.genPrefix(initPrefix)

	return fmt.Sprintf("%s%0*d", prefix, cg.NumLength, initNum)
}

func (cg *CodeGen) GenTotalCode(total int) []string {
	cg.mu.Lock()
	var codeIndex = cg.InitCode
	// 测试协程调度-共享内存并发竞争数据不一致
	//runtime.Gosched()
	cg.InitCode += total
	cg.mu.Unlock()
	data := make([]string, 0, total)
	initPrefix := codeIndex / cg.numWeight
	initNum := codeIndex % cg.numWeight
	prefix := cg.genPrefix(initPrefix)
	for i := 0; i < total; i++ {
		code := fmt.Sprintf("%s%0*d", prefix, cg.NumLength, initNum)
		data = append(data, code)
		if initNum == 999 {
			initPrefix += 1
			initNum = 0
			prefix = cg.genPrefix(initPrefix)
		} else {
			initNum++
		}
	}
	return data
}

func (cg *CodeGen) GenCodeNum(code string) int {
	ret := 0
	runes := []rune(code)
	for _, c := range runes {
		ret = 26*ret + (int(c-'A') + 1)
	}
	return ret
}

func (cg *CodeGen) initNumWeight() {
	weight := 1
	for i := 0; i < cg.NumLength; i++ {
		weight *= 10
	}
	cg.numWeight = weight
}
