package goo_kafka

import (
	"fmt"
	goo_redis "github.com/liqiongtao/googo.io/goo-redis"
	"testing"
	"time"
)

var (
	topic   = "test-01"
	groupId = "test-01"
)

func TestProducer(t *testing.T) {
	redis, _ := goo_redis.New(goo_redis.Config{
		Addr:     "redis.in:20063",
		Password: "",
		DB:       0,
	})

	Init(Config{
		User:     "admin",
		Password: "",
		Addrs:    []string{"kafka.in:20092"},
		Redis:    redis,
	})

	for i := 0; i < 20; i++ {
		go Producer().WithKey("test:5").SendMessage(topic, []byte(fmt.Sprintf("%d", i)))
	}

	time.Sleep(3 * time.Second)
}

func TestConsumer(t *testing.T) {
	redis, _ := goo_redis.New(goo_redis.Config{
		Addr:     "redis.in:20063",
		Password: "",
		DB:       0,
	})

	Init(Config{
		User:              "admin",
		Password:          "",
		Addrs:             []string{"kafka.in:20092"},
		HeartbeatInterval: 10,
		SessionTimeout:    30,
		RebalanceTimeout:  45,
		Redis:             redis,
	})

	Consumer().ConsumeGroup(groupId, []string{topic}, func(msg *ConsumerMessage, consumerErr *ConsumerError) error {
		fmt.Println(time.Now().Format("15:04:05"), string(msg.Value))

		time.Sleep(5 * time.Second)

		return nil
	})
}

func TestConsumer2(t *testing.T) {
	Init(Config{
		User:              "admin",
		Password:          "",
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
		Password:          "",
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
