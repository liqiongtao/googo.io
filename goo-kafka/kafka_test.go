package gookafka

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	gooredis "github.com/liqiongtao/googo.io/goo-redis"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
)

var (
	topic   = "test-01"
	groupId = "test-01"
)

type TestMessage struct {
	Id      int    `json:"id"`
	TraceId string `json:"trace-id"`
}

func (t *TestMessage) Topic() string {
	return topic
}

func (t *TestMessage) Key() string {
	return fmt.Sprintf("%s:%d", topic, t.Id)
}

func (t *TestMessage) Headers() map[string]string {
	return map[string]string{
		"source": "my-test",
	}
}

func (t *TestMessage) Serialize() []byte {
	t.TraceId = goo_utils.UUID()
	b, _ := json.Marshal(t)
	return b
}

func (t *TestMessage) Deserialize(b []byte) {
	if err := json.Unmarshal(b, t); err != nil {
		fmt.Println(err)
	}
}

func TestProducer(t *testing.T) {
	redis, _ := gooredis.New(gooredis.Config{
		Addr:     "redis.in:20063",
		Password: "",
		DB:       0,
	})
	Init(Config{
		User:     "admin",
		Password: "",
		Addrs:    []string{"kafka.in:20092"},
	}, RedisOption(redis))

	for i := 0; i < 20; i++ {
		Producer().SendMessage(&TestMessage{Id: 200 + i})
	}

	time.Sleep(3 * time.Second)
}

func TestConsumer(t *testing.T) {
	Init(Config{
		User:     "admin",
		Password: "",
		Addrs:    []string{"kafka.in:20092"},
		RedisConfig: gooredis.Config{
			Addr:     "redis.in:20063",
			Password: "",
			DB:       0,
		},
	})

	Consumer().ConsumeGroup(groupId, []string{topic}, func(ctx context.Context, msg *ConsumerMessage, consumerErr *ConsumerError) error {
		time.Sleep(5 * time.Second)

		return nil
	})
}
