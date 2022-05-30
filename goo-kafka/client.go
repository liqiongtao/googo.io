package goo_kafka

import (
	"github.com/Shopify/sarama"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"time"
)

type client struct {
	user     string
	password string
	addrs    []string
	timeout  time.Duration
	sarama.Client
}

func (cli *client) init() (err error) {
	config := sarama.NewConfig()

	if cli.user != "" {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = cli.user
		config.Net.SASL.Password = cli.password
	}

	// 等所有follower都成功后再返回
	config.Producer.RequiredAcks = sarama.WaitForAll
	// 分区策略为Hash，解决相同key的消息落在一个分区
	config.Producer.Partitioner = sarama.NewHashPartitioner
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.AutoCommit.Enable = true              // 自动提交
	config.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second // 间隔
	config.Consumer.Offsets.Retry.Max = 3
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky

	config.ChannelBufferSize = 1000
	config.Version = sarama.V0_10_2_0

	if cli.timeout > 0 {
		config.Producer.Timeout = cli.timeout
	}

	cli.Client, err = sarama.NewClient(cli.addrs, config)
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
