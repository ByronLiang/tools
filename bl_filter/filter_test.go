package bl_filter

import (
	"testing"
)

func TestNewFilter(t *testing.T) {
	storage := NewQueue(3)
	f := NewFilter(storage)
	f.Add("2:ss")
	f.Add("1:zz")
	res := f.All()
	t.Log(res)
	// 更新当前元素值
	f.UpdateCurrent("3:zzz")
	el := f.Current()
	t.Log(el)
	// 新增元素
	f.Add("4:kk")
	t.Log("before max: ", f.All())
	// 触发元素回收, 移除旧元素
	f.Add("5:cc")
	t.Log("after max: ", f.All())
}
