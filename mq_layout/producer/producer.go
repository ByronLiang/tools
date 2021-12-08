package producer

import (
	"fmt"
	"time"
)

type producer struct {
	config         BufferConfig
	buffer         chan ContentBody
	isSetRateLimit bool
	isStop         bool
	rate           *time.Ticker
	errChan        chan error
}

func NewProducer(config BufferConfig) *producer {
	if config.FrequencyLimit == 0 {
		return &producer{config: config}
	}
	rate := time.NewTicker(time.Duration(config.FrequencyLimit) * time.Second)
	// 开启一个协程 定时从管道取出数据进行发送
	p := &producer{
		config:         config,
		buffer:         make(chan ContentBody, config.Size),
		isSetRateLimit: true,
		isStop:         false,
		rate:           rate,
		errChan:        make(chan error, config.Size),
	}
	go func(p *producer) {
		for {
			select {
			case <-p.rate.C:
				if p.isStop {
					return
				}
				p.releaseBuffer()
			}
		}
	}(p)
	return p
}

func (p *producer) releaseBuffer() {
release:
	for i := 0; i < p.config.SendSize; i++ {
		select {
		case body, ok := <-p.buffer:
			if ok {
				if err := body.Send(); err != nil {
					p.errChan <- err
				}
			}
		default:
			break release
		}
	}
}

func (p *producer) SendBuffer(body ContentBody) (bool, error) {
	if !p.isSetRateLimit {
		return false, fmt.Errorf("producer no set rate limit")
	}
	if p.isStop {
		return false, fmt.Errorf("producer had stop")
	}
	select {
	case p.buffer <- body:
		return true, nil
	default:
		return false, nil
	}
}

func (p *producer) Send(body ContentBody) error {
	return body.Send()
}

func (p *producer) Stop() error {
	if !p.isSetRateLimit {
		return nil
	}
	p.isStop = true
	p.rate.Stop()
	time.Sleep(500 * time.Millisecond)
	close(p.buffer)
	// 将chan 剩余进行消费
	for body := range p.buffer {
		if err := body.Send(); err != nil {
			p.errChan <- err
		}
	}
	time.Sleep(500 * time.Millisecond)
	close(p.errChan)
	return nil
}
