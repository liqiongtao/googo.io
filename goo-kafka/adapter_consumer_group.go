package goo_kafka

import (
	"fmt"
	"github.com/IBM/sarama"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"time"
)

// 分组
type group struct {
	*client
	id      string
	handler ConsumerHandler
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
			if key == "" {
				key = g.GetKey(claim.Topic(), string(msg.Value))
			}
			log.WithField("key", key)

			// 建立缓存
			if g.redis != nil {
				if g.redis.Exists(key).Val() > 0 {
					session.MarkMessage(msg, "")
					continue
				}
				g.redis.Set(key, time.Now().Format("2006-01-02 15:04:05"), time.Hour)
			}

			if err := g.handler(&ConsumerMessage{ConsumerMessage: msg, GroupSession: session}, nil); err != nil {
				log.Error(err)
				continue
			}

			// 删除缓存
			if g.redis != nil {
				g.redis.Del(key)
			}

			// 提交
			session.MarkMessage(msg, "")
		}
	}
}
