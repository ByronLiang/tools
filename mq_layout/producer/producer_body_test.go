package producer

import "fmt"

type MqMessage struct {
	Key  string
	Data string
}

func (m *MqMessage) GetTopic() string {
	return "tracking"
}

func (m *MqMessage) GetContent() string {
	return m.Data
}

func (m *MqMessage) Send() error {
	// 使用mq client 发送到broker 并确保mq 生产者不丢失消息
	fmt.Println("send to kafka mq", m.GetContent())
	return nil
}
