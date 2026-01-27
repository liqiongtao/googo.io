package goo_kafka

import (
	goo_redis "github.com/liqiongtao/googo.io/goo-redis"
)

var (
	__client *Client
)

// 初始化
func Init(conf Config, opts ...Option) error {
	c, err := New(conf, opts...)
	if err != nil {
		return err
	}
	__client = c
	return nil
}

// 初始化
func New(conf Config, opts ...Option) (*Client, error) {
	c := &Client{conf: conf}
	for _, opt := range opts {
		switch opt.Name {
		case RedisName:
			c.redis = opt.Value.(*goo_redis.Client)
		}
	}
	if err := c.init(); err != nil {
		return nil, err
	}
	return c, nil
}

// 客户端
func DefClient() *Client {
	return __client
}

// 生产者
func Producer(opts ...Option) IProducer {
	return __client.Producer()
}

// 消费者
func Consumer() IConsumer {
	return __client.Consumer()
}

// 主题列表
func Topics() []string {
	if __client == nil {
		return []string{}
	}
	return __client.Topics()
}

// 分区数量
func Partitions(topic string) []int32 {
	if __client == nil {
		return []int32{}
	}
	return __client.Partitions(topic)
}

// 分区数量
func OffsetInfo(topic, groupId string) (data []map[string]int64) {
	if __client == nil {
		return []map[string]int64{}
	}
	return __client.OffsetInfo(topic, groupId)
}
