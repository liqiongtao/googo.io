package goo_kafka

import (
	"github.com/Shopify/sarama"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"time"
)

type client struct {
	addrs   []string
	timeout time.Duration
	sarama.Client
}

func (cli *client) init() (err error) {
	config := sarama.NewConfig()

	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky

	config.ChannelBufferSize = 1000
	config.Version = sarama.V0_10_2_0

	if cli.timeout > 0 {
		config.Producer.Timeout = cli.timeout
	}

	cli.Client, err = sarama.NewClient(cli.addrs, config)
	if err != nil {
		goo_log.Error(err)
	}

	return
}

func (cli *client) Close() {
	if !cli.Client.Closed() {
		cli.Client.Close()
	}
}
