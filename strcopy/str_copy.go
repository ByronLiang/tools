package strcopy

import "unsafe"

type Build struct {
	Prefix    string
	container []byte
	copyIndex int
}

func NewBuild(prefix string) *Build {
	container := make([]byte, 128)
	copy(container[:len(prefix)], prefix)
	return &Build{
		Prefix:    prefix,
		container: container,
		copyIndex: len(prefix),
	}
}

func (b *Build) Write(str string) string {
	copy(b.container[b.copyIndex:], str)
	return *(*string)(unsafe.Pointer(&b.container))
}
