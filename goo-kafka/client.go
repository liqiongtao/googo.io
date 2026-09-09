package goo_kafka

import (
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	goo_redis "github.com/liqiongtao/googo.io/goo-redis"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
)

type Client struct {
	conf Config
	sarama.Client
	redis *goo_redis.Client
}

func (c *Client) init() (err error) {
	id := uuid.New().String()
	config := sarama.NewConfig()

	if c.conf.User != "" {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = c.conf.User
		config.Net.SASL.Password = c.conf.Password
	}

	config.ClientID = id
	config.ChannelBufferSize = 1024
	config.Version = sarama.V3_0_0_0

	// 等所有follower都成功后再返回
	config.Producer.RequiredAcks = sarama.WaitForAll
	// 分区策略为Manual，指定分区发送消息
	//config.Producer.Partitioner = sarama.NewManualPartitioner
	// 分区策略为Hash，解决相同key的消息落在一个分区
	//config.Producer.Partitioner = sarama.NewHashPartitioner
	// 分区策略为Random，解决消费组分布式部署
	config.Producer.Partitioner = sarama.NewRandomPartitioner
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
	config.Consumer.Group.InstanceId = id

	c.Client, err = sarama.NewClient(c.conf.Addrs, config)
	if err != nil {
		goo_log.WithTag("goo-kafka").Error(err)
	}

	if cfg := c.conf.RedisConfig; cfg.Addr != "" {
		var redisErr error
		c.redis, redisErr = goo_redis.New(cfg)
		if redisErr != nil {
			goo_log.WithTag("goo-kafka").Error("Redis 初始化失败", redisErr)
			c.redis = nil
			if err == nil {
				err = redisErr
			}
		}
	}

	return
}

func (c *Client) Close() {
	if !c.Client.Closed() {
		c.Client.Close()
	}
}

func (c *Client) GetKey(topic, msg string) string {
	return fmt.Sprintf("goo:mq:%s:%s", time.Now().Format("20060102"), goo_utils.MD5([]byte(topic+msg)))
}

func (c *Client) Redis() *goo_redis.Client {
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
		goo_log.WithTag("goo-kafka").Error(err)
		return []string{}
	}

	return topics
}

// 分区数量
func (c *Client) Partitions(topic string) []int32 {
	partitions, err := c.Client.Partitions(topic)
	if err != nil {
		goo_log.WithTag("goo-kafka").WithField("topic", topic).Error(err)
		return []int32{}
	}

	return partitions
}

// 分区数量
func (c *Client) OffsetInfo(topic, groupId string) (data []map[string]int64) {
	data = []map[string]int64{}

	partitions := Partitions(topic)
	if l := len(partitions); l == 0 {
		return
	}

	var (
		l = goo_log.WithTag("goo-kafka").WithField("groupId", groupId).WithField("topic", topic)
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

		nextOffset, msg := pom.NextOffset()
		if msg != "" {
			l.Error(msg)
			continue
		}

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
