package goo_kafka

import (
	goo_context "github.com/liqiongtao/googo.io/goo-context"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
)

var (
	__client *client
)

func Init(conf Config) error {
	__client = &client{conf: conf}
	goo_utils.AsyncFunc(func() {
		select {
		case <-goo_context.Cancel().Done():
			__client.Close()
			return
		}
	})
	return __client.init()
}

func Client() *client {
	return __client
}

func Producer() IProducer {
	return &producer{client: __client}
}

func Consumer() IConsumer {
	return &consumer{client: __client}
}
