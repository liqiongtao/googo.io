package goo_kafka

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
	"github.com/liqiongtao/googo.io/goocontext"
)

var errConcurrentConsume = errors.New("concurrent consume")

// 分组
type group struct {
	cli     *Client
	id      string
	handler ConsumerHandler
}

func (g group) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (group) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (g group) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case <-session.Context().Done():
			// rebalance / shutdown 属正常结束，勿当错误触发外层重试噪音
			return nil

		case msg, ok := <-claim.Messages():
			if !ok {
				return nil
			}
			// 同一条处理到成功，或判定为需停 claim 的失败（避免 Mark 跨过）
			for {
				var err error
				func() {
					defer func() {
						if r := recover(); r != nil {
							goo_log.WithTag("goo-kafka-consumer-group", g.id).Error(r)
							err = fmt.Errorf("panic: %v", r)
						}
					}()
					err = g.doHandler(msg, session)
				}()
				if err == nil {
					break
				}
				// 处理失败不退出 claim（否则整 session 结束触发 rebalance）；等 1s 再试本条
				select {
				case <-session.Context().Done():
					return nil
				case <-time.After(time.Second):
				}
			}
		}
	}
}

func (g group) doHandler(msg *sarama.ConsumerMessage, session sarama.ConsumerGroupSession) (err error) {
	key := string(msg.Key)

	m := goo_utils.M{
		"topic":     msg.Topic,
		"key":       key,
		"partition": msg.Partition,
		"offset":    msg.Offset,
		"timestamp": msg.Timestamp.Format("2006-01-02 15:04:05"),
	}

	{
		if len(msg.Value) > 0 {
			var body any
			if err = json.Unmarshal(msg.Value, &body); err == nil {
				m["body"] = body
			} else {
				m["body"] = string(msg.Value)
			}
		}

		if len(msg.Headers) > 0 {
			headers := map[string]string{}
			for _, i := range msg.Headers {
				headers[string(i.Key)] = string(i.Value)
			}
			m["headers"] = headers
		}
	}

	// 跟随 session：rebalance 时可取消；不挂 Root，进程退出只停拉取
	ctx := goocontext.WithGenerateTraceId(session.Context())
	log := goocontext.Log(ctx).WithTag("goo-kafka-consumer-group", g.id).WithField("msg", m)

	{
		var uniqKey string
		if key != "" {
			uniqKey = consumeLockKey(g.id, msg.Topic, key)
		} else {
			// 无 Key 时用分区+offset，避免相同 payload 误互斥
			uniqKey = consumeLockKey(g.id, msg.Topic, fmt.Sprintf("%d-%d", msg.Partition, msg.Offset))
		}
		if g.cli.redis != nil {
			ok, setErr := g.cli.redis.SetNX(uniqKey, goo_utils.M{
				"topic":     msg.Topic,
				"body":      m["body"],
				"headers":   m["headers"],
				"timestamp": m["timestamp"],
			}.String(), 300*time.Second).Result()
			if setErr != nil {
				// Redis 故障时不能当去重命中，否则会 MarkMessage 丢消息
				log.Error("消息去重失败", setErr)
				err = setErr
				return
			}
			if !ok {
				log.Warn("消息并发消费中，稍后重试")
				err = errConcurrentConsume
				return
			}
			defer func() {
				g.cli.redis.Del(uniqKey)
			}()
		}
	}

	if g.cli.redis != nil && key != "" {
		g.cli.redis.Set(consumeCacheKey(msg.Topic, key), goo_utils.M{
			"topic":     msg.Topic,
			"body":      m["body"],
			"headers":   m["headers"],
			"timestamp": m["timestamp"],
		}.String(), time.Hour)
	}

	t1 := time.Now()
	defer func() {
		if g.cli.redis != nil && key != "" {
			g.cli.redis.Expire(consumeCacheKey(msg.Topic, key), 5*time.Second)
		}

		log = log.WithField("执行时间", fmt.Sprintf("%f", float64(time.Now().Sub(t1).Milliseconds())/1e3))
		if err != nil {
			log.Error("消息消费失败", err)
			return
		}

		log.Debug("消息消费成功")
	}()

	if err = g.handler(ctx, &ConsumerMessage{ConsumerMessage: msg, GroupSession: session}, nil); err != nil {
		return
	}

	session.MarkMessage(msg, "")

	return
}
