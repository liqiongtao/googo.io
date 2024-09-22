package goo_kafka

import (
	"fmt"
	"github.com/IBM/sarama"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
	"os"
	"strconv"
	"time"
)

type client struct {
	conf Config
	sarama.Client
}

func (c client) init() (err error) {
	id := strconv.Itoa(os.Getpid())
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

	return
}

func (c *client) Close() {
	if !c.Client.Closed() {
		c.Client.Close()
	}
}

func (c *client) GetKey(topic, msg string) string {
	return fmt.Sprintf("goo:mq:%s:%s", time.Now().Format("20060102"), goo_utils.MD5([]byte(topic+msg)))
}
