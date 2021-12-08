package consumer

import (
	"fmt"
	"sync"
)

type Consumer struct {
	BufferSet map[string]*Buffer
	mu        sync.RWMutex
}

func (c *Consumer) AddBuffer(buffer *Buffer) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.BufferSet[buffer.Topic] = buffer
	return
}

func (c *Consumer) GetBuffer(topic string) (*Buffer, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if buffer, ok := c.BufferSet[topic]; ok {
		return buffer, nil
	}
	return nil, fmt.Errorf("topic no exist: %s", topic)
}

func (c *Consumer) DelBuffer(topic string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.BufferSet[topic]; ok {
		delete(c.BufferSet, topic)
		return true
	}
	return false
}

func NewConsumer(buffer ...*Buffer) *Consumer {
	c := &Consumer{
		BufferSet: make(map[string]*Buffer),
		mu:        sync.RWMutex{},
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, b := range buffer {
		c.BufferSet[b.Topic] = b
	}
	return c
}
