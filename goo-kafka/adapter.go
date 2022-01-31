package goo_kafka

// 生产者
type iProducer interface {
	init() error

	Close()

	// 发送消息 - 同步
	SendMessage(topic string, message []byte) (partition int32, offset int64, err error)

	// 发送消息 - 异步
	SendAsyncMessage(topic string, message []byte) (partition int32, offset int64, err error)
}

// 消费者
type iConsumer interface {
	init() error

	Close()

	// 处理分区消息
	// partition: 分区ID
	// offset: 消息偏移量，-2=从头开始，-1=获取最新的
	PartitionConsume(topic string, partition int32, offset int64, handler ConsumerHandler)

	// 处理指定偏移量消息
	// partition: 分区ID，默认为0
	// offset: 消息偏移量，-2=从头开始，-1=获取最新的
	Consume(topic string, offset int64, handler ConsumerHandler)

	// 处理最新消息
	// partition: 分区ID，默认为0
	// offset: 消息偏移量，-2=从头开始，-1=获取最新的
	ConsumeNewest(topic string, handler ConsumerHandler)

	// 处理所有消息，从第一条开始
	// partition: 分区ID，默认为0
	// offset: 消息偏移量，-2=从头开始，-1=获取最新的
	ConsumeOldest(topic string, handler ConsumerHandler)

	// 分组topic
	ConsumeGroup(groupId string, topics []string, handler ConsumerHandler)
}
