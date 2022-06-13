package goo_kafka

import (
	"github.com/Shopify/sarama"
)

type ConsumerHandler func(msg *ConsumerMessage, consumerErr *ConsumerError) error

type ConsumerMessage struct {
	*sarama.ConsumerMessage
}

type ConsumerError struct {
	*sarama.ConsumerError
}
