package common

import (
	"bytes"
	"encoding/gob"
)

/**
深拷贝工具方法

若拷贝对象是struct, 成员有Slice/Map数据类型, 浅拷贝下, 对这类成员进行数据变更, 会影响原对象

若使用深拷贝, 拷贝对象与被拷贝对象不会互相影响; 拷贝的对象中没有引用类型，只需浅拷贝即可
*/
func DeepCopy(dst, src interface{}) error {
	defer func() {

	}()
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}
