package goo_kafka

import (
	"github.com/Shopify/sarama"
	goo_context "github.com/liqiongtao/googo.io/goo-context"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
)

type consumer struct {
	*client

	hasSetPartition bool  // 是否设置分区
	partition       int32 // 分区

	offset int64
}

func (c *consumer) Client() sarama.Client {
	return c.client.Client
}

// 设置 分区
func (c *consumer) WithPartition(partition int32) iConsumer {
	c.hasSetPartition = true
	c.partition = partition
	return c
}

// 设置 起始位置
func (c *consumer) WithOffset(offset int64) iConsumer {
	c.offset = offset
	return c
}

// 设置 起始位置 = 最新位置
func (c *consumer) WithOffsetNewest() iConsumer {
	c.offset = sarama.OffsetNewest
	return c
}

// 设置 起始位置 = 从头开始
func (c *consumer) WithOffsetOldest() iConsumer {
	c.offset = sarama.OffsetOldest
	return c
}

// 消费消息，默认处理最新消息
func (c *consumer) Consume(topic string, handler ConsumerHandler) {
	consumer, err := sarama.NewConsumerFromClient(c.Client())
	if err != nil {
		goo_log.WithTag("goo-kafka-consumer").Error(err)
		return
	}
	defer consumer.Close()

	if c.offset == 0 {
		c.offset = sarama.OffsetNewest
	}

	pc, err := consumer.ConsumePartition(topic, c.partition, c.offset)
	if err != nil {
		goo_log.WithTag("goo-kafka-consumer").Error(err)
		return
	}
	defer pc.Close()

	for {
		select {
		case <-goo_context.Cancel().Done():
			return

		case msg := <-pc.Messages():
			handler(&ConsumerMessage{msg}, nil)

		case err := <-pc.Errors():
			handler(nil, &ConsumerError{err})
		}
	}
}

// 分组
func (c *consumer) ConsumeGroup(groupId string, topics []string, handler ConsumerHandler) {
	cg, err := sarama.NewConsumerGroupFromClient(groupId, c.Client())
	if err != nil {
		goo_log.WithTag("goo-kafka-consumer-group").Error(err)
		return
	}

	g := group{handler: handler}

	goo_utils.AsyncFunc(func() {
		defer cg.Close()

		for {
			select {
			case <-goo_context.Cancel().Done():
				return

			case err := <-cg.Errors():
				goo_log.WithTag("goo-kafka-consumer-group").Error(err)
			}
		}
	})

	if err := cg.Consume(goo_context.Cancel(), topics, g); err != nil {
		goo_log.WithTag("goo-kafka-consumer-group").Error(err)
	}
}
