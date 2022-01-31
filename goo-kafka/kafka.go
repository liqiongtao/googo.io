package goo_kafka

import (
	"time"
)

var (
	__client *client
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
	return &producer{client: __client}
}

func Consumer() iConsumer {
	return &consumer{client: __client}
}
