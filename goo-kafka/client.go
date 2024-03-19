package goo_kafka

import (
	"github.com/IBM/sarama"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"os"
	"strconv"
	"time"
)

type client struct {
	conf Config
	sarama.Client
}

func (cli *client) init() (err error) {
	id := strconv.Itoa(os.Getpid())
	config := sarama.NewConfig()

	if cli.conf.User != "" {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = cli.conf.User
		config.Net.SASL.Password = cli.conf.Password
	}

	config.ClientID = id
	config.ChannelBufferSize = 1024
	config.Version = sarama.V3_0_0_0

	// 等所有follower都成功后再返回
	config.Producer.RequiredAcks = sarama.WaitForAll
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
		sarama.NewBalanceStrategySticky(),
		sarama.NewBalanceStrategyRange(),
	}
	config.Consumer.Group.Heartbeat.Interval = 10 * time.Second
	config.Consumer.Group.Session.Timeout = 30 * time.Second
	config.Consumer.Group.InstanceId = id

	cli.Client, err = sarama.NewClient(cli.conf.Addrs, config)
	if err != nil {
		goo_log.WithTag("goo-kafka").Error(err)
	}

	return
}

func (cli *client) Close() {
	if !cli.Client.Closed() {
		cli.Client.Close()
	}
}
