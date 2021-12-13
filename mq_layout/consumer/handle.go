package consumer

import (
	"context"
	"log"
	"time"
)

type GroupHandle interface {
	SetBuffer(buffer *Buffer) // 缓冲对象设置
	GetBuffer() *Buffer
	GetConsumeFrequency() int                    // 配置缓冲区消费频率
	ConsumeBufferHandle(data interface{})        // 消费缓冲区数据业务实现
	InitConsumerGroup(ctx context.Context) error // 初始化消费组
}

func DefaultConsumerHandle(buffer *Buffer) error {
	ctx, cancel := context.WithCancel(context.Background())
	buffer.groupHandle.SetBuffer(buffer)
	initConsumerGroupErrChan := make(chan error)
	// 从kafka 取出消息
	go func() {
		// 依赖外部初始化消费组
		err := buffer.groupHandle.InitConsumerGroup(ctx)
		if err != nil {
			// 初始化异常, 结束消费线程
			initConsumerGroupErrChan <- err
		}
	}()
	ticker := time.NewTicker(time.Duration(buffer.groupHandle.GetConsumeFrequency()) * time.Second)
	for {
		select {
		case <-ticker.C:
		loadData:
			for i := 0; i < buffer.Size; i++ {
				select {
				case res := <-buffer.Data:
					buffer.groupHandle.ConsumeBufferHandle(res)
				default:
					break loadData
				}
			}
		case initConsumerGroupErr := <-initConsumerGroupErrChan: // 消费组初始化异常
			ticker.Stop()
			close(initConsumerGroupErrChan)
			return initConsumerGroupErr
			// 从chan 消费缓冲内容
		case <-buffer.CloseSig:
			log.Println("close sig")
			// 监听关闭缓冲区
			ticker.Stop()
			time.Sleep(500 * time.Millisecond)
			// 关闭管道
			close(buffer.Data)
			// 将剩余缓冲区消息进行消费
			lastLen := len(buffer.Data)
			for i := 0; i < lastLen; i++ {
				select {
				case res := <-buffer.Data:
					buffer.groupHandle.ConsumeBufferHandle(res)
				default:
					continue
				}
			}
			time.Sleep(500 * time.Millisecond)
			// 关闭消费组成员
			cancel()
			return nil
		}
	}
}
