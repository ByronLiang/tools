package producer

import (
	"fmt"
	"testing"
	"time"
)

func TestNewProducer(t *testing.T) {
	c := BufferConfig{
		Size:           5,
		SendSize:       3,
		FrequencyLimit: 2,
	}
	p := NewProducer(c)
	for i := 0; i < 20; i++ {
		m := &MqMessage{
			Key:  "test",
			Data: fmt.Sprintf("mq-test-%d", i),
		}
		res, err := p.SendBuffer(m)
		if err == nil && res == false {
			// retry process
			for j := 1; j <= 3; j++ {
				time.Sleep(500 * 1 << j * time.Millisecond)
				res, _ := p.SendBuffer(m)
				if res {
					break
				}
			}
		}
		time.Sleep(300 * time.Millisecond)
	}
	err := p.Stop()
	if err != nil {
		t.Log(err)
		return
	}
	_, err = p.SendBuffer(&MqMessage{
		Key:  "test-stop",
		Data: "test-stop",
	})
	if err != nil {
		t.Log(err)
		return
	}
}
