package goo_kafka

import (
	"github.com/Shopify/sarama"
	goo_context "github.com/liqiongtao/googo.io/goo-context"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
)

type consumer struct {
	*client
}

// 处理分区消息
func (c *consumer) PartitionConsume(topic string, partition int32, offset int64, handler ConsumerHandler) {
	consumer, err := sarama.NewConsumerFromClient(c.Client)
	if err != nil {
		goo_log.WithTrace().Error(err)
		return
	}
	defer consumer.Close()

	pc, err := consumer.ConsumePartition(topic, partition, offset)
	if err != nil {
		goo_log.WithTrace().Error(err)
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

	return
}

// 处理指定偏移量消息
func (c *consumer) Consume(topic string, offset int64, handler ConsumerHandler) {
	c.PartitionConsume(topic, 0, offset, handler)
}

// 处理最新消息
func (c *consumer) ConsumeNewest(topic string, handler ConsumerHandler) {
	c.PartitionConsume(topic, 0, sarama.OffsetNewest, handler)
}

// 处理所有消息，从第一条开始
func (c *consumer) ConsumeOldest(topic string, handler ConsumerHandler) {
	c.PartitionConsume(topic, 0, sarama.OffsetOldest, handler)
}

// 处理分区消息
func (c *consumer) ConsumeGroup(groupId string, topics []string, handler ConsumerHandler) {
	cg, err := sarama.NewConsumerGroupFromClient(groupId, c.Client)
	if err != nil {
		goo_log.WithTrace().Error(err)
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
				goo_log.WithTrace().Error(err)
			}
		}
	})

	if err := cg.Consume(goo_context.Cancel(), topics, g); err != nil {
		goo_log.WithTrace().Error(err)
		return
	}
}
