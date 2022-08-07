# kafka

```
goo_kafka.Client().Topics()

// 发送消息，不指定分区
goo_kafka.Producer().SendMessage("test", []byte("hi hnatao"))

// 发送消息，指定分区
goo_kafka.Producer().WithPartition(0).SendMessage("test", []byte("hi hnatao"))

// 发送异步消息，不指定分区
goo_kafka.Producer().SendAsyncMessage("test", []byte("hi hnatao"), func(msg *goo_kafka.ProducerMessage, err error) {
})

// 发送异步消息，指定分区
goo_kafka.Producer().WithPartition(0).SendAsyncMessage("test", []byte("hi hnatao"), func(msg *goo_kafka.ProducerMessage, err error) {
})

// 消费消息，指定分区，指定起始位置
goo_kafka.Consumer().WithPartition(0).WithOffset(100).Consume("test", func(msg *goo_kafka.ConsumerMessage, consumerErr *goo_kafka.ConsumerError) error {
    return nil
})

// 消费消息，指定分区，从最新位置开始
goo_kafka.Consumer().WithPartition(0).WithOffsetNewest().Consume("test", func(msg *goo_kafka.ConsumerMessage, consumerErr *goo_kafka.ConsumerError) error {
    return nil
})

// 消费消息，指定分区，从最头开始
goo_kafka.Consumer().WithPartition(0).WithOffsetOldest().Consume("test", func(msg *goo_kafka.ConsumerMessage, consumerErr *goo_kafka.ConsumerError) error {
    return nil
})

// 消费消息，分组消息，分组里面只要1个消费者消费
goo_kafka.Consumer().ConsumeGroup("test-id", []string{"test"}, func(msg *goo_kafka.ConsumerMessage, consumerErr *goo_kafka.ConsumerError) error {
    return nil
})
```