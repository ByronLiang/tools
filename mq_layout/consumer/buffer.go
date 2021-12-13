package consumer

type Buffer struct {
	Topic       string
	GroupId     string
	Data        chan interface{}
	Size        int
	isStop      bool
	CloseSig    chan struct{}
	handle      BufferHandle
	groupHandle GroupHandle
}

type BufferHandle func(buffer *Buffer) error

func NewBuffer(topic, groupId string, size int, handle BufferHandle) *Buffer {
	buffer := &Buffer{
		Topic:    topic,
		GroupId:  groupId,
		Data:     make(chan interface{}, size),
		Size:     size,
		CloseSig: make(chan struct{}, 1),
		handle:   handle,
	}
	return buffer
}

func NewGroupBuffer(topic, groupId string, size int, groupHandle GroupHandle) *Buffer {
	buffer := &Buffer{
		Topic:       topic,
		GroupId:     groupId,
		Data:        make(chan interface{}, size),
		Size:        size,
		CloseSig:    make(chan struct{}, 1),
		groupHandle: groupHandle,
	}
	return buffer
}

func (b *Buffer) AddData(data interface{}) bool {
	select {
	case b.Data <- data:
		return true
	default:
		return false
	}
}

func (b *Buffer) GetDataSize() int {
	return len(b.Data)
}

func (b *Buffer) EmptyBuffer() {
	b.Data = make(chan interface{}, b.Size)
}

func (b *Buffer) GetAllData() []interface{} {
	res := make([]interface{}, 0)
	for data := range b.Data {
		res = append(res, data)
	}
	return res
}

func (b *Buffer) IsStop() bool {
	return b.isStop
}

func (b *Buffer) SetStop() {
	b.isStop = true
}

func (b *Buffer) Close() bool {
	if b.IsStop() {
		return true
	}
	select {
	case b.CloseSig <- struct{}{}:
		b.isStop = true
		return true
	default:
		return false
	}
}
