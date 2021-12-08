package consumer

type Buffer struct {
	Topic  string
	Data   chan interface{}
	Size   int
	isStop bool
	close  chan struct{}
}

type BufferHandle func(buffer *Buffer)

func NewBuffer(topic string, size int, handle BufferHandle) *Buffer {
	buffer := &Buffer{
		Topic: topic,
		Data:  make(chan interface{}, size),
		Size:  size,
		close: make(chan struct{}, 1),
	}
	go func(buffer *Buffer) {
		handle(buffer)
	}(buffer)
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

func (b *Buffer) Close() bool {
	select {
	case b.close <- struct{}{}:
		b.isStop = true
		return true
	default:
		return false
	}
}
