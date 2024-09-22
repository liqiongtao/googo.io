package goo_kafka

import "github.com/IBM/sarama"

// 生产者
type IProducer interface {
	init() error

	Close()

	Client() sarama.Client

	// 发送消息到指定分区
	WithPartition(partition int32) IProducer

	// 指定Key
	WithKey(key string) IProducer

	// 发送消息 - 同步
	SendMessage(topic string, message []byte) (partition int32, offset int64, err error)

	// 发送消息 - 异步
	SendAsyncMessage(topic string, message []byte, cb MessageHandler) (err error)
}

// 消费者
type IConsumer interface {
	init() error

	Close()

	Client() sarama.Client

	// 从指定分区消费
	WithPartition(partition int32) IConsumer

	// 从指定位置开始
	WithOffset(offset int64) IConsumer

	// 从最新位置开始
	WithOffsetNewest() IConsumer

	// 从头开始
	WithOffsetOldest() IConsumer

	// 消费
	Consume(topic string, handler ConsumerHandler)

	// 分组topic
	ConsumeGroup(groupId string, topics []string, handler ConsumerHandler)
}
