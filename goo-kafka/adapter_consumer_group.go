package goo_kafka

import (
	"fmt"
	"github.com/IBM/sarama"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
)

// 分组
type group struct {
	id      string
	handler ConsumerHandler
	config  Config
}

func (g group) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (group) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (g group) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	log := goo_log.WithTag("goo-kafka-consumer-group").
		WithField("groupId", g.id).
		WithField("topic", claim.Topic()).
		WithField("partition", claim.Partition())

	for {
		select {
		case <-session.Context().Done():
			return fmt.Errorf("关闭会话上下文: %s", session.Context().Err())

		case msg, ok := <-claim.Messages():
			if !ok {
				return fmt.Errorf("消费通道关闭: groupId=%s topic=%s partition=%d", g.id, claim.Topic(), claim.Partition())
			}

			key := string(msg.Key)
			log.WithField("key", key)

			if err := g.handler(&ConsumerMessage{ConsumerMessage: msg, GroupSession: session}, nil); err != nil {
				log.Error(err)
				continue
			}

			// 删除缓存
			if redis := g.config.Redis; redis != nil && key != "" {
				redis.Del(key)
			}

			session.MarkMessage(msg, "")
		}
	}
}
