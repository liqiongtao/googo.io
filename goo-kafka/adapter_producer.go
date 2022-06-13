package goo_kafka

import (
	"github.com/Shopify/sarama"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
)

type producer struct {
	*client
}

// 发送消息 - 同步
func (p *producer) SendMessage(topic string, message []byte) (partition int32, offset int64, err error) {
	var producer sarama.SyncProducer

	producer, err = sarama.NewSyncProducerFromClient(p.Client)
	if err != nil {
		goo_log.WithTag("goo-kafka-producer").Error(err)
		return
	}
	defer producer.Close()

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(message),
		Key:   sarama.StringEncoder(topic),
	}

	return producer.SendMessage(msg)
}

// 发送消息 - 异步
func (p *producer) SendAsyncMessage(topic string, message []byte) (partition int32, offset int64, err error) {
	var producer sarama.AsyncProducer

	producer, err = sarama.NewAsyncProducerFromClient(p.Client)
	if err != nil {
		goo_log.WithTag("goo-kafka-producer").Error(err)
		return
	}
	defer producer.Close()

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(message),
		Key:   sarama.StringEncoder(topic),
	}

	producer.Input() <- msg

	select {
	case msg := <-producer.Successes():
		partition = msg.Partition
		offset = msg.Offset

	case pe := <-producer.Errors():
		partition = pe.Msg.Partition
		offset = pe.Msg.Offset
		err = pe.Err
		goo_log.WithTag("goo-kafka-producer").Error(err.Error())
	}

	return
}
