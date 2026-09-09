package goo_kafka

import (
	"sync"

	goo_redis "github.com/liqiongtao/googo.io/goo-redis"
)

var (
	__client *Client
	__mu     sync.RWMutex
)

// 初始化
func Init(conf Config, opts ...Option) error {
	c, err := New(conf, opts...)
	if err != nil {
		return err
	}
	__mu.Lock()
	__client = c
	__mu.Unlock()
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
	__mu.RLock()
	defer __mu.RUnlock()
	return __client
}

// 生产者
func Producer(opts ...Option) IProducer {
	__mu.RLock()
	cli := __client
	__mu.RUnlock()
	if cli == nil {
		return nil
	}
	return cli.Producer(opts...)
}

// 消费者
func Consumer() IConsumer {
	__mu.RLock()
	cli := __client
	__mu.RUnlock()
	if cli == nil {
		return nil
	}
	return cli.Consumer()
}

// 主题列表
func Topics() []string {
	__mu.RLock()
	cli := __client
	__mu.RUnlock()
	if cli == nil {
		return []string{}
	}
	return cli.Topics()
}

// 分区数量
func Partitions(topic string) []int32 {
	__mu.RLock()
	cli := __client
	__mu.RUnlock()
	if cli == nil {
		return []int32{}
	}
	return cli.Partitions(topic)
}

// 分区数量
func OffsetInfo(topic, groupId string) (data []map[string]int64) {
	__mu.RLock()
	cli := __client
	__mu.RUnlock()
	if cli == nil {
		return []map[string]int64{}
	}
	return cli.OffsetInfo(topic, groupId)
}
