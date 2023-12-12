package goo_kafka

import (
	"github.com/IBM/sarama"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
)

type producer struct {
	*client

	msg *sarama.ProducerMessage
}

func (p *producer) Client() sarama.Client {
	return p.client.Client
}

// 指定分区
func (p *producer) WithPartition(partition int32) iProducer {
	p.msg.Partition = partition
	return p
}

// 发送消息 - 同步
func (p *producer) SendMessage(topic string, message []byte) (partition int32, offset int64, err error) {
	p.msg.Topic = topic
	p.msg.Value = sarama.ByteEncoder(message)
	p.msg.Key = sarama.StringEncoder(topic)

	var producer sarama.SyncProducer

	producer, err = sarama.NewSyncProducerFromClient(p.Client())
	if err != nil {
		goo_log.WithTag("goo-kafka-producer").
			WithField("topic", topic).
			WithField("msg", string(message)).
			Error(err)
		return
	}
	defer producer.Close()

	return producer.SendMessage(p.msg)
}

// 发送消息 - 异步
func (p *producer) SendAsyncMessage(topic string, message []byte, cb MessageHandler) (err error) {
	p.msg.Topic = topic
	p.msg.Value = sarama.ByteEncoder(message)
	p.msg.Key = sarama.StringEncoder(topic)

	var producer sarama.AsyncProducer

	producer, err = sarama.NewAsyncProducerFromClient(p.Client())
	if err != nil {
		goo_log.WithTag("goo-kafka-producer").
			WithField("topic", topic).
			WithField("msg", string(message)).
			Error(err)
		return
	}
	defer producer.Close()

	producer.Input() <- p.msg

	select {
	case msg := <-producer.Successes():
		cb(&ProducerMessage{msg}, nil)
	case e := <-producer.Errors():
		goo_log.WithTag("goo-kafka-producer").
			WithField("topic", topic).
			WithField("msg", string(message)).
			Error(e.Msg)
		cb(&ProducerMessage{e.Msg}, e.Err)
	}

	return
}
