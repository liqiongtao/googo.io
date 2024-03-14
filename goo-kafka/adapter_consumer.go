package goo_kafka

import (
	"errors"
	"github.com/IBM/sarama"
	goo_context "github.com/liqiongtao/googo.io/goo-context"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
	"time"
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
	l := goo_log.WithTag("goo-kafka-consumer").WithField("topic", topic)

	consumer, err := sarama.NewConsumerFromClient(c.Client())
	if err != nil {
		l.Error(err)
		return
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			l.Error(err)
		}
	}()

	if c.offset == 0 {
		c.offset = sarama.OffsetNewest
	}

	pc, err := consumer.ConsumePartition(topic, c.partition, c.offset)
	if err != nil {
		l.Error(err)
		return
	}
	defer func() {
		if err := pc.Close(); err != nil {
			l.Error(err)
		}
	}()

	for {
		select {
		case <-goo_context.Cancel().Done():
			l.Debug("Context被取消,停止消费")
			return

		case err := <-pc.Errors():
			if err != nil {
				l.Error(err)
			}

		case msg, ok := <-pc.Messages():
			if !ok {
				l.Debug("消息通道被关闭,停止消费")
				return
			}
			handler(&ConsumerMessage{msg}, nil)
		}
	}
}

// 分组
func (c *consumer) ConsumeGroup(groupId string, topics []string, handler ConsumerHandler) {
	l := goo_log.WithTag("goo-kafka-consumer-group").
		WithField("groupId", groupId).
		WithField("topics", topics)

	cg, err := sarama.NewConsumerGroupFromClient(groupId, c.Client())
	if err != nil {
		l.Error(err)
		return
	}
	defer func() {
		if err := cg.Close(); err != nil {
			l.Error(err)
		}
	}()

	goo_utils.AsyncFunc(func() {
		for {
			select {
			case <-goo_context.Cancel().Done():
				if err := cg.Close(); err != nil {
					l.Error(err)
				}
				return

			case err := <-cg.Errors():
				if err != nil {
					l.Error(err)
				}
			}
		}
	})

	g := group{id: groupId, handler: handler}
	if err := cg.Consume(goo_context.Cancel(), topics, g); err != nil && !errors.Is(err, sarama.ErrClosedConsumerGroup) {
		l.Error(err)
	}

	time.Sleep(time.Second)

	l.Debug("订阅消息退出")
}
