package gookafka

import (
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	goolog "github.com/liqiongtao/googo.io/goo-log"
	gooredis "github.com/liqiongtao/googo.io/goo-redis"
)

type Client struct {
	conf Config
	sarama.Client
	redis *gooredis.Client
}

func (c *Client) init() (err error) {
	clientID := c.conf.ClientID
	if clientID == "" {
		clientID = uuid.New().String()
	}
	config := sarama.NewConfig()

	if c.conf.User != "" {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = c.conf.User
		config.Net.SASL.Password = c.conf.Password
	}

	config.ClientID = clientID
	config.ChannelBufferSize = 1024
	config.Version = sarama.V3_0_0_0

	// 等所有follower都成功后再返回
	config.Producer.RequiredAcks = sarama.WaitForAll
	// 分区策略为 Hash，相同 key 落入同一分区，保证同 key 有序
	config.Producer.Partitioner = sarama.NewHashPartitioner
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.Retry.Max = 5
	config.Producer.Timeout = 10 * time.Second

	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.AutoCommit.Enable = true          // 自动提交
	config.Consumer.Offsets.AutoCommit.Interval = time.Second // 间隔
	config.Consumer.Offsets.Retry.Max = 5
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{
		sarama.NewBalanceStrategyRoundRobin(),
		sarama.NewBalanceStrategySticky(),
		sarama.NewBalanceStrategyRange(),
	}
	config.Consumer.Group.Heartbeat.Interval = 10 * time.Second
	// session.timeout = heartbeat.interval * 4
	config.Consumer.Group.Session.Timeout = 60 * time.Second
	// rebalance.timeout = session.timeout * 1.5
	config.Consumer.Group.Rebalance.Timeout = 90 * time.Second
	if c.conf.HeartbeatInterval > 0 {
		config.Consumer.Group.Heartbeat.Interval = time.Duration(c.conf.HeartbeatInterval) * time.Second
	}
	if c.conf.SessionTimeout > 0 {
		config.Consumer.Group.Session.Timeout = time.Duration(c.conf.SessionTimeout) * time.Second
	}
	if c.conf.RebalanceTimeout > 0 {
		config.Consumer.Group.Rebalance.Timeout = time.Duration(c.conf.RebalanceTimeout) * time.Second
	}
	// 仅在配置了稳定 InstanceId 时启用静态成员，避免每次随机 ID 加重 rebalance
	if c.conf.InstanceId != "" {
		config.Consumer.Group.InstanceId = c.conf.InstanceId
	}

	c.Client, err = sarama.NewClient(c.conf.Addrs, config)
	if err != nil {
		goolog.WithTag("goo-kafka").Error(err)
		return
	}

	if c.redis == nil {
		if cfg := c.conf.RedisConfig; cfg.Addr != "" {
			var redisErr error
			c.redis, redisErr = gooredis.New(cfg)
			if redisErr != nil {
				goolog.WithTag("goo-kafka").Error("Redis 初始化失败", redisErr)
				c.redis = nil
				_ = c.Client.Close()
				c.Client = nil
				err = redisErr
			}
		}
	}

	return
}

func (c *Client) Close() {
	if c == nil || c.Client == nil {
		return
	}
	if !c.Client.Closed() {
		_ = c.Client.Close()
	}
}

func (c *Client) GetKey(topic, key string) string {
	return prodDedupKey(topic, key)
}

// prodDedupKey 生产端去重 key，与消费端隔离
func prodDedupKey(topic, key string) string {
	return fmt.Sprintf("goo:kafka:prod:%s:%s", topic, key)
}

// consumeCacheKey 消费端缓存 key，与生产端隔离
func consumeCacheKey(topic, key string) string {
	return fmt.Sprintf("goo:kafka:cache:%s:%s", topic, key)
}

// consumeLockKey 消费组并发锁，含 group/topic，带命名空间
func consumeLockKey(groupId, topic, key string) string {
	return fmt.Sprintf("goo:kafka:lock:%s:%s:%s", groupId, topic, key)
}

func (c *Client) Redis() *gooredis.Client {
	return c.redis
}

// 生产者
func (c *Client) Producer(opts ...Option) IProducer {
	var focus bool
	for _, opt := range opts {
		switch opt.Name {
		case FocusName:
			focus = opt.Value.(bool)
		}
	}
	return &producer{cli: c, focus: focus}
}

// 消费者
func (c *Client) Consumer() IConsumer {
	return &consumer{cli: c}
}

// 题列表
func (c *Client) Topics() []string {
	topics, err := c.Client.Topics()
	if err != nil {
		goolog.WithTag("goo-kafka").Error(err)
		return []string{}
	}

	return topics
}

// 分区数量
func (c *Client) Partitions(topic string) []int32 {
	partitions, err := c.Client.Partitions(topic)
	if err != nil {
		goolog.WithTag("goo-kafka").WithField("topic", topic).Error(err)
		return []int32{}
	}

	return partitions
}

// 分区数量
func (c *Client) OffsetInfo(topic, groupId string) (data []map[string]int64) {
	data = []map[string]int64{}

	partitions := c.Partitions(topic)
	if l := len(partitions); l == 0 {
		return
	}

	var (
		l = goolog.WithTag("goo-kafka").WithField("groupId", groupId).WithField("topic", topic)
	)

	om, err := sarama.NewOffsetManagerFromClient(groupId, c.Client)
	if err != nil {
		l.Error(err)
		return
	}
	defer om.Close()

	for _, partition := range partitions {
		offset, err := c.GetOffset(topic, partition, -1)
		if err != nil {
			l.Error(err)
			continue
		}

		pom, err := om.ManagePartition(topic, partition)
		if err != nil {
			l.Error(err)
			continue
		}

		nextOffset, _ := pom.NextOffset()

		backlog := offset
		if nextOffset != -1 {
			backlog -= nextOffset
		}

		data = append(data, map[string]int64{
			"partition":  int64(partition),
			"offset":     offset,
			"nextOffset": nextOffset,
			"backlog":    backlog,
		})
	}

	return
}
