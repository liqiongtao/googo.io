package goo_kafka

import (
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
		User:     "",
		Password: "",
		Addrs:    []string{"127.0.0.1:9002"},
	})

	//for i := 0; i < 8; i++ {
	//	Producer().SendMessage(topic, []byte(fmt.Sprintf("1-%d", i)))
	//}

	Producer().SendMessage(topic, []byte("2-2"))
}

func TestConsumer(t *testing.T) {
	Init(Config{
		User:              "",
		Password:          "",
		Addrs:             []string{"127.0.0.1:9002"},
		HeartbeatInterval: 1,
		SessionTimeout:    10,
		RebalanceTimeout:  5,
	})

	Consumer().ConsumeGroup(groupId, []string{topic}, func(msg *ConsumerMessage, consumerErr *ConsumerError) error {
		fmt.Println(time.Now().Format("15:04:05"), string(msg.Value), msg.Timestamp, msg.BlockTimestamp)

		time.Sleep(60 * time.Second)

		switch string(msg.Value) {
		case "1-1", "1-4", "1-7", "1-8", "1-5":
			msg.Commit()
		}

		return nil
	})
}

func TestConsumer2(t *testing.T) {
	Init(Config{
		User:              "",
		Password:          "",
		Addrs:             []string{"127.0.0.1:9002"},
		HeartbeatInterval: 1,
		SessionTimeout:    12,
		RebalanceTimeout:  15,
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
		User:              "",
		Password:          "",
		Addrs:             []string{"127.0.0.1:9002"},
		HeartbeatInterval: 1,
		SessionTimeout:    15,
		RebalanceTimeout:  15,
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
