package goo_kafka

import (
	"github.com/Shopify/sarama"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
)

type consumer struct {
	*client
}

// 处理分区消息
func (c *consumer) PartitionConsume(topic string, partition int32, offset int64) (pc sarama.PartitionConsumer, err error) {
	var consumer sarama.Consumer

	consumer, err = sarama.NewConsumerFromClient(c.c)
	if err != nil {
		goo_log.Error(err)
		return
	}

	pc, err = consumer.ConsumePartition(topic, partition, offset)
	if err != nil {
		goo_log.Error(err)
	}

	return
}

// 处理指定偏移量消息
func (c *consumer) Consume(topic string, offset int64) (pc sarama.PartitionConsumer, err error) {
	return c.PartitionConsume(topic, 0, offset)
}

// 处理最新消息
func (c *consumer) ConsumeNewest(topic string) (pc sarama.PartitionConsumer, err error) {
	return c.PartitionConsume(topic, 0, sarama.OffsetNewest)
}

// 处理所有消息，从第一条开始
func (c *consumer) ConsumeOldest(topic string) (pc sarama.PartitionConsumer, err error) {
	return c.PartitionConsume(topic, 0, sarama.OffsetOldest)
}
