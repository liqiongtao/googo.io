# kafka

```
gookafka.Client().Topics()

// 发送消息，不指定分区
gookafka.Producer().SendMessage("test", []byte("hi hnatao"))

// 发送消息，指定分区
gookafka.Producer().WithPartition(0).SendMessage("test", []byte("hi hnatao"))

// 发送异步消息，不指定分区
gookafka.Producer().SendAsyncMessage("test", []byte("hi hnatao"), func(msg *gookafka.ProducerMessage, err error) {
})

// 发送异步消息，指定分区
gookafka.Producer().WithPartition(0).SendAsyncMessage("test", []byte("hi hnatao"), func(msg *gookafka.ProducerMessage, err error) {
})

// 消费消息，指定分区，指定起始位置
gookafka.Consumer().WithPartition(0).WithOffset(100).Consume("test", func(msg *gookafka.ConsumerMessage, consumerErr *gookafka.ConsumerError) error {
    return nil
})

// 消费消息，指定分区，从最新位置开始
gookafka.Consumer().WithPartition(0).WithOffsetNewest().Consume("test", func(msg *gookafka.ConsumerMessage, consumerErr *gookafka.ConsumerError) error {
    return nil
})

// 消费消息，指定分区，从最头开始
gookafka.Consumer().WithPartition(0).WithOffsetOldest().Consume("test", func(msg *gookafka.ConsumerMessage, consumerErr *gookafka.ConsumerError) error {
    return nil
})

// 消费消息，分组消息，分组里面只要1个消费者消费
gookafka.Consumer().ConsumeGroup("test-id", []string{"test"}, func(msg *gookafka.ConsumerMessage, consumerErr *gookafka.ConsumerError) error {
    return nil
})
```