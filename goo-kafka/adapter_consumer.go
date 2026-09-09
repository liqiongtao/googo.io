package goo_kafka

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/IBM/sarama"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
	"github.com/liqiongtao/googo.io/goocontext"
)

type consumer struct {
	cli *Client

	hasSetPartition bool  // 是否设置分区
	partition       int32 // 分区

	hasSetOffset bool  // 是否设置起始位置
	offset       int64 // 起始位置（0 为合法绝对 offset）
}

func (c *consumer) Client() sarama.Client {
	return c.cli.Client
}

// 设置 分区
func (c *consumer) WithPartition(partition int32) IConsumer {
	c.hasSetPartition = true
	c.partition = partition
	return c
}

// 设置 起始位置
func (c *consumer) WithOffset(offset int64) IConsumer {
	c.hasSetOffset = true
	c.offset = offset
	return c
}

// 设置 起始位置 = 最新位置
func (c *consumer) WithOffsetNewest() IConsumer {
	c.hasSetOffset = true
	c.offset = sarama.OffsetNewest
	return c
}

// 设置 起始位置 = 从头开始
func (c *consumer) WithOffsetOldest() IConsumer {
	c.hasSetOffset = true
	c.offset = sarama.OffsetOldest
	return c
}

// 消费消息，默认处理最新消息
func (c *consumer) Consume(topic string, handler ConsumerHandler) {
	log := goo_log.WithTag("goo-kafka-consumer").WithField("topic", topic)

	consumer, err := sarama.NewConsumerFromClient(c.Client())
	if err != nil {
		log.Error(err)
		return
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Error(err)
		}
	}()

	offset := c.offset
	if !c.hasSetOffset {
		offset = sarama.OffsetNewest
	}

	partitions := []int32{c.partition}
	if !c.hasSetPartition {
		partitions = c.cli.Partitions(topic)
		if len(partitions) == 0 {
			log.Error("no partitions")
			return
		}
	}

	var wg sync.WaitGroup
	for _, partition := range partitions {
		pc, err := consumer.ConsumePartition(topic, partition, offset)
		if err != nil {
			log.WithField("partition", partition).Error(err)
			continue
		}
		wg.Add(1)
		go func(pc sarama.PartitionConsumer, partition int32) {
			defer wg.Done()
			defer func() {
				if err := pc.Close(); err != nil {
					log.WithField("partition", partition).Error(err)
				}
			}()
			c.consumePartition(topic, pc, handler, log.WithField("partition", partition))
		}(pc, partition)
	}
	wg.Wait()
}

func (c *consumer) consumePartition(topic string, pc sarama.PartitionConsumer, handler ConsumerHandler, log *goo_log.Entry) {
	for {
		select {
		case <-goocontext.Root().Done():
			log.Debug("Context被取消,停止消费")
			return

		case err, ok := <-pc.Errors():
			if !ok {
				log.Debug("错误通道被关闭,停止消费")
				return
			}
			if err != nil {
				log.Error(err)
			}

		case msg, ok := <-pc.Messages():
			if !ok {
				log.Debug("消息通道被关闭,停止消费")
				return
			}

			ctx := goocontext.WithGenerateTraceId(context.Background())
			func() {
				defer goo_utils.Recovery()
				if err := handler(ctx, &ConsumerMessage{ConsumerMessage: msg}, nil); err != nil {
					log.Error(err)
				}
			}()

			key := string(msg.Key)
			if c.cli.redis != nil {
				c.cli.redis.Del(key)
			}
		}
	}
}

// 分组
func (c *consumer) ConsumeGroup(groupId string, topics []string, handler ConsumerHandler) {
	l := goo_log.WithTag("goo-kafka-consumer-group").
		WithField("groupId", groupId).
		WithField("topics", topics)

	cg, err := sarama.NewConsumerGroupFromClient(groupId, c.cli.Client)
	if err != nil {
		l.Error(err)
		return
	}
	defer func() {
		cg.Close()
		l.Debug("consumer-group 退出")
	}()

	var (
		done = make(chan struct{})
	)

	ctx, cancel := context.WithCancel(goocontext.Root())
	defer cancel()

	goo_utils.AsyncFunc(func() {
		for err := range cg.Errors() {
			if err != nil {
				l.Error(err)
			}
		}
	})

	goo_utils.AsyncFunc(func() {
		defer close(done)
		for {
			if ctx.Err() != nil {
				return
			}
			err := cg.Consume(ctx, topics, group{id: groupId, handler: handler, cli: c.cli})
			if ctx.Err() != nil {
				return
			}
			if err != nil && !errors.Is(err, sarama.ErrClosedConsumerGroup) {
				l.Error(err)
				// 避免 broker 异常时空转打满 CPU
				select {
				case <-ctx.Done():
					return
				case <-time.After(time.Second):
				}
			}
		}
	})

	<-ctx.Done()
	cancel()
	<-done

	time.Sleep(time.Second)
}
