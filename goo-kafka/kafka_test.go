package goo_kafka

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

var (
	topic   = "test-01"
	groupId = "test-01"
)

func TestProducer(t *testing.T) {
	Init(Config{
		User:     "admin",
		Password: "7fdacd2183ab",
		Addrs:    []string{"kafka.in:20092"},
	})

	//for i := 0; i < 8; i++ {
	//	Producer().SendMessage(topic, []byte(fmt.Sprintf("1-%d", i)))
	//}

	Producer().SendMessage(topic, []byte("3"))
}

func TestConsumer(t *testing.T) {
	Init(Config{
		User:              "admin",
		Password:          "7fdacd2183ab",
		Addrs:             []string{"kafka.in:20092"},
		HeartbeatInterval: 10,
		SessionTimeout:    30,
		RebalanceTimeout:  45,
	})

	Consumer().ConsumeGroup(groupId, []string{topic}, func(msg *ConsumerMessage, consumerErr *ConsumerError) error {
		fmt.Println(time.Now().Format("15:04:05"), string(msg.Value))

		switch string(msg.Value) {
		case "1":
			return errors.New("异常")
		}

		return nil
	})
}

func TestConsumer2(t *testing.T) {
	Init(Config{
		User:              "admin",
		Password:          "7fdacd2183ab",
		Addrs:             []string{"kafka.in:20092"},
		HeartbeatInterval: 10,
		SessionTimeout:    30,
		RebalanceTimeout:  45,
	})

	Consumer().ConsumeGroup(groupId, []string{topic}, func(msg *ConsumerMessage, consumerErr *ConsumerError) error {
		fmt.Println(time.Now().Format("15:04:05"), string(msg.Value))
		fmt.Println(msg.Timestamp, msg.BlockTimestamp)

		time.Sleep(3 * time.Second)

		switch string(msg.Value) {
		case "1-1", "1-4", "1-7", "1-8":
			msg.Commit()
		}

		return nil
	})
}

func TestConsumer3(t *testing.T) {
	Init(Config{
		User:              "admin",
		Password:          "7fdacd2183ab",
		Addrs:             []string{"kafka.in:20092"},
		HeartbeatInterval: 10,
		SessionTimeout:    30,
		RebalanceTimeout:  45,
	})

	Consumer().ConsumeGroup(groupId, []string{topic}, func(msg *ConsumerMessage, consumerErr *ConsumerError) error {
		fmt.Println(time.Now().Format("15:04:05"), string(msg.Value))
		fmt.Println(msg.Timestamp, msg.BlockTimestamp)

		time.Sleep(3 * time.Second)

		switch string(msg.Value) {
		case "1-1", "1-4", "1-7", "1-8":
			msg.Commit()
		}

		return nil
	})
}
