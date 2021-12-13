package consumer

import (
	"fmt"
	"log"
	"sync"
)

type Consumer struct {
	BufferSet map[string]*Buffer
	mu        sync.RWMutex
	wg        sync.WaitGroup
}

func (c *Consumer) AddBuffer(buffer *Buffer) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if buffer.groupHandle != nil {
		c.wg.Add(1)
		c.BufferSet[buffer.Topic] = buffer
		go func(b *Buffer) {
			defer func() {
				log.Printf("topic: %s stop\n", b.Topic)
				c.wg.Done()
			}()
			err := DefaultConsumerHandle(b)
			if err != nil {
				b.SetStop()
				return
			}
		}(buffer)
		return
	}
	c.wg.Add(1)
	c.BufferSet[buffer.Topic] = buffer
	go func(b *Buffer) {
		defer func() {
			log.Printf("topic: %s stop\n", b.Topic)
			c.wg.Done()
		}()
		err := b.handle(b)
		if err != nil {
			b.SetStop()
			return
		}
	}(buffer)
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
		buffer, _ := c.GetBuffer(topic)
		buffer.Close()
		delete(c.BufferSet, topic)
		return true
	}
	return false
}

func (c *Consumer) Close() {
	for topic, buffer := range c.BufferSet {
		ok := buffer.Close()
		if !ok {
			log.Printf("close buffer error: %s", topic)
		}
	}
	c.wg.Wait()
}

func NewConsumer(buffers ...*Buffer) *Consumer {
	c := &Consumer{
		BufferSet: make(map[string]*Buffer),
		mu:        sync.RWMutex{},
		wg:        sync.WaitGroup{},
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, buffer := range buffers {
		if buffer.groupHandle != nil {
			c.wg.Add(1)
			c.BufferSet[buffer.Topic] = buffer
			go func(b *Buffer) {
				defer func() {
					log.Printf("topic: %s stop\n", b.Topic)
					c.wg.Done()
				}()
				err := DefaultConsumerHandle(b)
				if err != nil {
					b.SetStop()
					return
				}
			}(buffer)
			continue
		}
		c.wg.Add(1)
		c.BufferSet[buffer.Topic] = buffer
		go func(b *Buffer) {
			defer func() {
				c.wg.Done()
			}()
			err := b.handle(b)
			if err != nil {
				b.SetStop()
				return
			}
		}(buffer)
	}
	return c
}
