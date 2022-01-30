package goo_kafka

import (
	"time"
)

var (
	__client   *client
	__producer iProducer
	__consumer iConsumer
)

func Init(addrs ...string) error {
	__client = &client{
		addrs:   addrs,
		timeout: 5 * time.Second,
	}
	return __client.init()
}

func Client() *client {
	return __client
}

func Producer() iProducer {
	if __producer == nil {
		__producer = &producer{client: __client}
	}
	return __producer
}

func Consumer() iConsumer {
	if __consumer == nil {
		__consumer = &consumer{client: __client}
	}
	return __consumer
}
