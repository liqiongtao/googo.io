package goo_kafka

import (
	"errors"
	"github.com/IBM/sarama"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
	"time"
)

type producer struct {
	*client
	focus bool // 是否强制发送
}

func (p *producer) Client() sarama.Client {
	return p.client.Client
}

// 发送消息 - 同步
func (p *producer) SendMessage(msg IMessage) (partition int32, offset int64, err error) {
	m := &sarama.ProducerMessage{
		Topic: msg.Topic(),
		Value: sarama.ByteEncoder(msg.Serialize()),
	}

	if v := msg.Key(); v != "" {
		m.Key = sarama.StringEncoder(v)
	}
	if data := msg.Headers(); data != nil {
		var headers []sarama.RecordHeader
		for k, v := range data {
			headers = append(headers, sarama.RecordHeader{
				Key:   []byte(k),
				Value: []byte(v),
			})
		}
		m.Headers = headers
	}

	defer func() {
		log := goo_log.WithTag("goo-kafka-producer").WithField("msg", goo_utils.M{
			"topic":     msg.Topic(),
			"key":       msg.Key(),
			"headers":   msg.Headers(),
			"body":      msg,
			"partition": m.Partition,
			"offset":    m.Offset,
		})
		if err != nil {
			log.Error("消息发送失败", err)
			return
		}
		log.Debug("消息发送成功")
	}()

	// 添加缓存
	if p.redis != nil && len(msg.Key()) > 0 {
		if p.focus {
			p.redis.Del(msg.Key())
		}
		if p.redis.Exists(msg.Key()).Val() > 0 {
			err = errors.New("KEY已存在")
			return
		}
		p.redis.Set(msg.Key(), goo_utils.M{
			"topic":     msg.Topic(),
			"body":      msg,
			"headers":   msg.Headers(),
			"timestamp": time.Now().Format("2006-01-02 15:04:05"),
		}.String(), time.Hour)
	}

	var producer sarama.SyncProducer

	producer, err = sarama.NewSyncProducerFromClient(p.Client())
	if err != nil {
		return
	}
	defer producer.Close()

	partition, offset, err = producer.SendMessage(m)

	return
}

// 发送消息 - 异步
func (p *producer) SendAsyncMessage(msg IMessage, cb MessageHandler) (err error) {
	m := &sarama.ProducerMessage{
		Topic: msg.Topic(),
		Value: sarama.ByteEncoder(msg.Serialize()),
	}

	if v := msg.Key(); v != "" {
		m.Key = sarama.StringEncoder(v)
	}
	if data := msg.Headers(); data != nil {
		var headers []sarama.RecordHeader
		for k, v := range data {
			headers = append(headers, sarama.RecordHeader{
				Key:   []byte(k),
				Value: []byte(v),
			})
		}
		m.Headers = headers
	}

	defer func() {
		log := goo_log.WithTag("goo-kafka-producer").WithField("msg", goo_utils.M{
			"topic":     msg.Topic(),
			"key":       msg.Key(),
			"headers":   msg.Headers(),
			"body":      msg,
			"partition": m.Partition,
			"offset":    m.Offset,
		})
		if err != nil {
			log.Error("消息发送失败", err)
			return
		}
		log.Debug("消息发送成功")
	}()

	// 添加缓存
	if p.redis != nil && len(msg.Key()) > 0 {
		if p.focus {
			p.redis.Del(msg.Key())
		}
		if p.redis.Exists(msg.Key()).Val() > 0 {
			err = errors.New("KEY已存在")
			return
		}
		p.redis.Set(msg.Key(), goo_utils.M{
			"topic":     msg.Topic(),
			"body":      msg,
			"headers":   msg.Headers(),
			"timestamp": time.Now().Format("2006-01-02 15:04:05"),
		}.String(), time.Hour)
	}

	var producer sarama.AsyncProducer

	producer, err = sarama.NewAsyncProducerFromClient(p.Client())
	if err != nil {
		return
	}
	defer producer.Close()

	producer.Input() <- m

	select {
	case msg := <-producer.Successes():
		cb(&ProducerMessage{msg}, nil)
	case e := <-producer.Errors():
		err = e.Err
		cb(&ProducerMessage{e.Msg}, e.Err)
	}

	return
}
