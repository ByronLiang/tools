package consumer

import (
	"fmt"
	"testing"
	"time"
)

func TestNewConsumer(t *testing.T) {
	eventBuffer := NewBuffer("event", "trace-event", 100, topicEventHandle)
	c := NewConsumer(eventBuffer)
	time.Sleep(5 * time.Second)
	buffer, err := c.GetBuffer("event")
	if err != nil {
		t.Fatal(err)
		return
	}
	buffer.Close()
	time.Sleep(2 * time.Second)
}

func topicEventHandle(b *Buffer) error {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			fmt.Println("topicEventHandle ticker")
			// 从chan 消费缓冲内容
			// 模拟从 mq 拉取消息放进缓冲区
			//case msg, ok := <-kafkaConsumerClient.Messages():
		case <-b.CloseSig:
			// 监听关闭缓冲区
			ticker.Stop()
			time.Sleep(500 * time.Millisecond)
			close(b.Data)
		}
	}
}
