package producer

type ContentBody interface {
	GetTopic() string
	GetContent() string
	Send() error
}
