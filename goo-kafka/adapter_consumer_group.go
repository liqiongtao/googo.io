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
		WithField("topic", claim.Topic())

	for {
		select {
		case <-sess.Context().Done():
			l.Debug("Context被取消,停止消费")
			return nil

		case msg, ok := <-claim.Messages():
			if !ok {
				l.Debug("消息通道被关闭,停止消费")
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
