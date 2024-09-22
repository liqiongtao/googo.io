package goo_kafka

import (
	"fmt"
	"github.com/IBM/sarama"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
	"time"
)

type producer struct {
	*client
	msg   *sarama.ProducerMessage
	focus bool // 是否强制发送
}

func (p *producer) Client() sarama.Client {
	return p.client.Client
}

// 指定分区
func (p *producer) WithPartition(partition int32) IProducer {
	p.Config().Producer.Partitioner = sarama.NewManualPartitioner
	p.msg.Partition = partition
	return p
}

// 指定Key
func (p *producer) WithKey(key string) IProducer {
	if key != "" {
		p.msg.Key = sarama.StringEncoder(key)
	}
	return p
}

// 发送消息 - 同步
func (p *producer) SendMessage(topic string, message []byte) (partition int32, offset int64, err error) {
	p.msg.Topic = topic
	p.msg.Value = sarama.ByteEncoder(message)
	if p.msg.Key == nil || p.msg.Key.Length() == 0 {
		key := fmt.Sprintf("goo:mq:%s:%s", time.Now().Format("20060102"), goo_utils.MD5([]byte(topic+string(message))))
		p.msg.Key = sarama.StringEncoder(key)
	}

	keyByte, _ := p.msg.Key.Encode()
	key := string(keyByte)

	log := goo_log.WithTag("goo-kafka-producer").
		WithField("key", key).
		WithField("topic", topic).
		WithField("msg", string(message))

	// 添加缓存
	if redis := p.client.conf.Redis; redis != nil {
		if p.focus {
			redis.Del(key)
		}
		if redis.Exists(key).Val() > 0 {
			log.Debug("消息发送失败，Key已存在")
			return
		}
		redis.Set(key, time.Now().Format("2006-01-02 15:04:05"), time.Hour)
	}

	var producer sarama.SyncProducer

	producer, err = sarama.NewSyncProducerFromClient(p.Client())
	if err != nil {
		log.Error(err)
		return
	}
	defer producer.Close()

	return producer.SendMessage(p.msg)
}

// 发送消息 - 异步
func (p *producer) SendAsyncMessage(topic string, message []byte, cb MessageHandler) (err error) {
	p.msg.Topic = topic
	p.msg.Value = sarama.ByteEncoder(message)
	if p.msg.Key == nil || p.msg.Key.Length() == 0 {
		key := fmt.Sprintf("goo:mq:%s:%s", time.Now().Format("20060102"), goo_utils.MD5([]byte(topic+string(message))))
		p.msg.Key = sarama.StringEncoder(key)
	}

	keyByte, _ := p.msg.Key.Encode()
	key := string(keyByte)

	log := goo_log.WithTag("goo-kafka-producer").
		WithField("key", key).
		WithField("topic", topic).
		WithField("msg", string(message))

	// 添加缓存
	if redis := p.client.conf.Redis; redis != nil {
		if p.focus {
			redis.Del(key)
		}
		if redis.Exists(key).Val() > 0 {
			log.Debug("消息发送失败，Key已存在")
			return
		}
		redis.Set(key, time.Now().Format("2006-01-02 15:04:05"), time.Hour)
	}

	var producer sarama.AsyncProducer

	producer, err = sarama.NewAsyncProducerFromClient(p.Client())
	if err != nil {
		log.Error(err)
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
