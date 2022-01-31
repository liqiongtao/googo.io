package goo_kafka

import (
	"github.com/Shopify/sarama"
)

// 分组
type group struct {
	handler ConsumerHandler
}

func (g group) Setup(sess sarama.ConsumerGroupSession) error {
	return nil
}

func (group) Cleanup(sess sarama.ConsumerGroupSession) error {
	return nil
}

func (g group) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		if err := g.handler(&ConsumerMessage{msg}, nil); err != nil {
			continue
		}
		//fmt.Printf("topic:%q partition:%d offset:%d  value:%s\n", msg.Topic, msg.Partition, msg.Offset, string(msg.Value))
		sess.MarkMessage(msg, "")
	}
	return nil
}
