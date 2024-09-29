package goo_kafka

import (
	"github.com/IBM/sarama"
	goo_context "github.com/liqiongtao/googo.io/goo-context"
)

type MessageHandler func(msg *ProducerMessage, err error)

type ProducerMessage struct {
	*sarama.ProducerMessage
}

type ConsumerHandler func(ctx *goo_context.Context, msg *ConsumerMessage, consumerErr *ConsumerError) error

type ConsumerMessage struct {
	*sarama.ConsumerMessage
	GroupSession sarama.ConsumerGroupSession
}

func (msg ConsumerMessage) Commit() {
	if msg.GroupSession == nil {
		return
	}
	msg.GroupSession.MarkMessage(msg.ConsumerMessage, "")
}

type ConsumerError struct {
	*sarama.ConsumerError
}
