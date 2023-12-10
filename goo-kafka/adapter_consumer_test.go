package goo_kafka

import (
	"fmt"
	"github.com/IBM/sarama"
	"testing"
)

func TestConsumer_Client(t *testing.T) {
	config := sarama.NewConfig()

	config.Net.SASL.Enable = true
	config.Net.SASL.User = "admin"
	config.Net.SASL.Password = "7fdacd2183ab"

	client, err := sarama.NewClient([]string{"kafka.in:20092"}, config)
	if err != nil {
		panic(err.Error())
	}

	var (
		topic = "smartcard-txt2json-v2-prod"
	)

	partitions, err := client.Partitions(topic)
	if err != nil {
		panic(err.Error())
	}

	om, err := sarama.NewOffsetManagerFromClient("smart-card-v2-prod", client)
	if err != nil {
		panic(err.Error())
	}

	var totalBacklog int64 = 0

	for _, partition := range partitions {
		offset, err := client.GetOffset(topic, partition, -1)
		if err != nil {
			panic(err.Error())
		}

		pom, err := om.ManagePartition(topic, partition)
		if err != nil {
			panic(err.Error())
		}

		var backlog int64

		n, string := pom.NextOffset()
		if string != "" {
			panic(string)
		}

		if n == -1 {
			backlog = offset
		} else {
			backlog = offset - n
		}

		totalBacklog += backlog
		fmt.Printf("partation %d, Kafka中下一条数据offset: %d, 消费者下次提交的offset: %d, 数据积压量: %d\n", partition, offset, n, backlog)
	}

	fmt.Printf("\n%s 的数据积压量总数为 %d\n\n", topic, totalBacklog)
}
