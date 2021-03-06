# 消息中间件模板

定义适合 Kafka 等MQ中间件的生产者与消费者 SDK 模型

## producer

- 消息实体需要实现`ContentBody`接口，通过外部控制生产者投递消息的行为

- 可配置缓冲区长度与投递频率, 对上游投递进行流量限制，防止消息中间件异常，积压消息；也可以直接投递到MQ

## consumer

- 配置读取消息缓冲区, 存放消息体，批量执行消息

- 配置不同的topic进行相应消费回调闭包函数

- 隔离不同topic的缓冲区长度与消息处理频率

### GroupHandle 配置

针对 Kafka 消费组模型的消息消费模板: 默认使用`DefaultConsumerHandle`

使用消费组需要实现接口方法，并且满足 kafka sdk 的消费组实现方法: `sarama.ConsumerGroupHandler`

```go
type GroupHandle interface {
	SetBuffer(buffer *Buffer)
	GetBuffer() *Buffer
	GetConsumeFrequency() int
	ConsumeBufferHandle(data interface{})
	InitConsumerGroup(ctx context.Context) error
}
```

### 自定义消费流程

自定义消息消费流程配置可直接配置使用`handle`成员, 可以在闭包函数拉取消息并放置缓冲区内，并且初始定时器, 定时消费缓冲区的数据
