package goo_kafka

import (
	"github.com/IBM/sarama"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
)

// 分组
type group struct {
	id      string
	handler ConsumerHandler
}

func (g group) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (group) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (g group) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	l := goo_log.WithTag("goo-kafka-consumer-group").
		WithField("groupId", g.id).
		WithField("topic", claim.Topic()).
		WithField("partition", claim.Partition())

	for {
		select {
		case <-sess.Context().Done():
			l.Debug("关闭会话上下文")
			return nil

		case msg, ok := <-claim.Messages():
			if !ok {
				l.Debug("消费通道关闭")
				return nil
			}

			if err := g.handler(&ConsumerMessage{msg}, nil); err != nil {
				l.Error(err)
				continue
			}

			sess.MarkMessage(msg, "")
		}
	}
}
