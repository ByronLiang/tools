package strcopy

type StrCp struct {
	Prefix     string
	containers map[int][]byte // support different length str appended prefix without applying memory
	copyIndex  int
}

func NewBuild(prefix string) *StrCp {
	container := make([]byte, 128)
	copy(container[:len(prefix)], prefix)
	return &StrCp{
		Prefix:    prefix,
		copyIndex: len(prefix),
	}
}

// only single process
func (s *StrCp) Write(str string) string {
	_, ok := s.containers[len(str)]
	if !ok {
		s.containers[len(str)] = make([]byte, s.copyIndex+len(str)) // apply memory for space
		copy(s.containers[len(str)][:s.copyIndex], s.Prefix)
	}
	copy(s.containers[len(str)][s.copyIndex:], str)
	return string(s.containers[len(str)])
}
