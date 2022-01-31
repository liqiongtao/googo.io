package goo_kafka

import (
	"github.com/Shopify/sarama"
)

type ConsumerHandler func(msg *ConsumerMessage, err *ConsumerError) error

type ConsumerMessage struct {
	*sarama.ConsumerMessage
}

type ConsumerError struct {
	*sarama.ConsumerError
}
