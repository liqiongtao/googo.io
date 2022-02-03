package goo_kafka

import (
	"encoding/json"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"testing"
)

func TestClient(t *testing.T) {
	Init("", "", "122.228.113.230:19092")

	goo_log.Debug(Client().Topics())
	goo_log.Debug(Client().Partitions("A101"))
}

func TestProducer_SendMessage(t *testing.T) {
	Init("", "", "122.228.113.230:19092")

	partition, offset, err := Producer().SendMessage("1", "A101", []byte("ok"))
	if err != nil {
		goo_log.Error(err)
		return
	}
	goo_log.Debug(partition, offset)
}

func TestProducer_SendAsyncMessage(t *testing.T) {
	Init("", "", "122.228.113.230:19092")

	partition, offset, err := Producer().SendAsyncMessage("1", "A101", []byte("ok"))
	if err != nil {
		goo_log.Error(err)
		return
	}
	goo_log.Debug(partition, offset)
}

func TestConsumer_PartitionConsume(t *testing.T) {
	Init("", "", "122.228.113.230:19092")

	Consumer().PartitionConsume("A101", 0, 0, func(msg *ConsumerMessage, err *ConsumerError) error {
		if err != nil {
			return err.Err
		}
		b, _ := json.Marshal(&msg)
		goo_log.Debug(string(b), string(msg.Value))
		return nil
	})
}

func TestConsumer_Consume(t *testing.T) {
	Init("", "", "122.228.113.230:19092")

	Consumer().Consume("A101", 3, func(msg *ConsumerMessage, err *ConsumerError) error {
		if err != nil {
			return err.Err
		}
		b, _ := json.Marshal(&msg)
		goo_log.Debug(string(b), string(msg.Value))
		return nil
	})
}

func TestConsumer_ConsumeNewest(t *testing.T) {
	Init("", "", "122.228.113.230:19092")

	Consumer().ConsumeNewest("A101", func(msg *ConsumerMessage, err *ConsumerError) error {
		if err != nil {
			return err.Err
		}
		b, _ := json.Marshal(&msg)
		goo_log.Debug(string(b), string(msg.Value))
		return nil
	})
}

func TestConsumer_ConsumeOldest(t *testing.T) {
	Init("", "", "122.228.113.230:19092")

	Consumer().ConsumeOldest("A101", func(msg *ConsumerMessage, err *ConsumerError) error {
		if err != nil {
			return err.Err
		}
		b, _ := json.Marshal(&msg)
		goo_log.Debug(string(b), string(msg.Value))
		return nil
	})
}

// 分组1
func TestConsumer_PartitionConsumeGroup(t *testing.T) {
	Init("", "", "122.228.113.230:19092")

	Consumer().ConsumeGroup("101", []string{"A101"}, func(msg *ConsumerMessage, _ *ConsumerError) error {
		b, _ := json.Marshal(msg)
		goo_log.Debug("--1--", string(b))
		return nil
	})
}

// 分组2
func TestConsumer_PartitionConsumeGroup2(t *testing.T) {
	Init("", "", "122.228.113.230:19092")

	Consumer().ConsumeGroup("101", []string{"A101"}, func(msg *ConsumerMessage, _ *ConsumerError) error {
		b, _ := json.Marshal(msg)
		goo_log.Debug("--2--", string(b))
		return nil
	})
}
