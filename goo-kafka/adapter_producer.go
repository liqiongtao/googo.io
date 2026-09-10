package goo_kafka

import (
	"errors"
	"time"

	"github.com/IBM/sarama"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
)

type producer struct {
	cli   *Client
	focus bool // 是否强制发送
}

func (p *producer) Client() sarama.Client {
	return p.cli.Client
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

	dedupKey, err := p.tryDedup(msg)
	if err != nil {
		return
	}

	var producer sarama.SyncProducer

	producer, err = sarama.NewSyncProducerFromClient(p.Client())
	if err != nil {
		p.clearDedup(dedupKey)
		return
	}
	defer producer.Close()

	partition, offset, err = producer.SendMessage(m)
	if err != nil {
		p.clearDedup(dedupKey)
	}

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

	dedupKey, err := p.tryDedup(msg)
	if err != nil {
		return
	}

	var producer sarama.AsyncProducer

	producer, err = sarama.NewAsyncProducerFromClient(p.Client())
	if err != nil {
		p.clearDedup(dedupKey)
		return
	}
	defer producer.Close()

	producer.Input() <- m

	select {
	case msg := <-producer.Successes():
		if cb != nil {
			cb(&ProducerMessage{msg}, nil)
		}
	case e := <-producer.Errors():
		err = e.Err
		p.clearDedup(dedupKey)
		if cb != nil {
			cb(&ProducerMessage{e.Msg}, e.Err)
		}
	}

	return
}

// tryDedup 用 SetNX 原子占位去重；focus 时用 Set 原子覆盖旧 key。
// 返回已占用的 dedupKey（空表示未启用去重）；发送失败时需 clearDedup。
func (p *producer) tryDedup(msg IMessage) (dedupKey string, err error) {
	if p.cli.redis == nil || len(msg.Key()) == 0 {
		return "", nil
	}

	dedupKey = prodDedupKey(msg.Topic(), msg.Key())
	val := goo_utils.M{
		"topic":     msg.Topic(),
		"body":      msg,
		"headers":   msg.Headers(),
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
	}.String()

	if p.focus {
		// Set 原子覆盖，避免 Del+SetNX 竞态导致双占位
		if setErr := p.cli.redis.Set(dedupKey, val, time.Hour).Err(); setErr != nil {
			return "", setErr
		}
		return dedupKey, nil
	}

	ok, setErr := p.cli.redis.SetNX(dedupKey, val, time.Hour).Result()
	if setErr != nil {
		return "", setErr
	}
	if !ok {
		return "", errors.New("KEY已存在")
	}
	return dedupKey, nil
}

func (p *producer) clearDedup(dedupKey string) {
	if dedupKey == "" || p.cli.redis == nil {
		return
	}
	p.cli.redis.Del(dedupKey)
}
