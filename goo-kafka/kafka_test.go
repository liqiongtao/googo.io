package goo_kafka

import (
	"encoding/json"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"testing"
)

func TestKafkaProducer_SendMessage(t *testing.T) {
	Init("122.228.113.230:19092")

	partition, offset, err := Producer().SendMessage("A101", []byte("ok"))
	if err != nil {
		goo_log.Error(err)
		return
	}

	goo_log.Debug(partition, offset)
}

func TestKafkaProducer_SendAsyncMessage(t *testing.T) {
	Init("122.228.113.230:19092")

	partition, offset, err := Producer().SendAsyncMessage("A101", []byte("ok"))
	if err != nil {
		goo_log.Error(err)
		return
	}

	goo_log.Debug(partition, offset)
}

func TestKafkaConsumer_PartitionConsume(t *testing.T) {
	Init("122.228.113.230:19092")

	pc, err := Consumer().PartitionConsume("A101", 0, 0)
	if err != nil {
		goo_log.Error(err)
		return
	}

	for {
		select {
		case msg := <-pc.Messages():
			b, _ := json.Marshal(&msg)
			goo_log.Debug(string(b), string(msg.Value))
		case err := <-pc.Errors():
			goo_log.Error(err)
		}
	}
}

func TestKafkaConsumer_Consume(t *testing.T) {
	Init("122.228.113.230:19092")

	pc, err := Consumer().Consume("A101", 3)
	if err != nil {
		goo_log.Error(err)
		return
	}

	for {
		select {
		case msg := <-pc.Messages():
			b, _ := json.Marshal(&msg)
			goo_log.Debug(string(b))
		case err := <-pc.Errors():
			goo_log.Error(err)
		}
	}
}

func TestKafkaConsumer_ConsumeNewest(t *testing.T) {
	Init("122.228.113.230:19092")

	pc, err := Consumer().ConsumeNewest("A101")
	if err != nil {
		goo_log.Error(err)
		return
	}

	for {
		select {
		case msg := <-pc.Messages():
			b, _ := json.Marshal(&msg)
			goo_log.Debug(string(b))
		case err := <-pc.Errors():
			goo_log.Error(err)
		}
	}
}

func TestKafkaConsumer_ConsumeOldest(t *testing.T) {
	Init("122.228.113.230:19092")

	pc, err := Consumer().ConsumeOldest("A101", )
	if err != nil {
		goo_log.Error(err)
		return
	}

	for {
		select {
		case msg := <-pc.Messages():
			b, _ := json.Marshal(&msg)
			goo_log.Debug(string(b))
		case err := <-pc.Errors():
			goo_log.Error(err)
		}
	}
}
