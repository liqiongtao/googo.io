package goo_kafka

import (
	goo_context "github.com/liqiongtao/googo.io/goo-context"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
	"time"
)

var (
	__client *client
)

func Init(user, password string, addrs ...string) error {
	__client = &client{
		user:     user,
		password: password,
		addrs:    addrs,
		timeout:  5 * time.Second,
	}
	goo_utils.AsyncFunc(func() {
		for {
			select {
			case <-goo_context.Cancel().Done():
				__client.Close()
				return
			}
		}
	})
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
