package goo_kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
	"github.com/liqiongtao/googo.io/goocontext"
)

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
			return fmt.Errorf("关闭会话上下文: %s", session.Context().Err())

		case msg, ok := <-claim.Messages():
			if !ok {
				return fmt.Errorf("消费通道关闭: groupId=%s topic=%s partition=%d", g.id, claim.Topic(), claim.Partition())
			}
			func() {
				defer goo_utils.Recovery()
				g.doHandler(msg, session)
			}()
		}
	}
}

func (g group) doHandler(msg *sarama.ConsumerMessage, session sarama.ConsumerGroupSession) (err error) {
	// 消息key
	key := string(msg.Key)

	// 消息试题
	m := goo_utils.M{
		"topic":     msg.Topic,
		"key":       key,
		"partition": msg.Partition,
		"offset":    msg.Offset,
		"timestamp": msg.Timestamp.Format("2006-01-02 15:04:05"),
	}

	// 填充数据
	{
		// body
		if len(msg.Value) > 0 {
			var body interface{}
			if err = json.Unmarshal(msg.Value, &body); err == nil {
				m["body"] = body
			} else {
				m["body"] = string(msg.Value)
			}
		}

		// headers
		if len(msg.Headers) > 0 {
			headers := map[string]string{}
			for _, i := range msg.Headers {
				headers[string(i.Key)] = string(i.Value)
			}
			m["headers"] = headers
		}
	}

	// 在途消费不挂 Root：进程退出只停拉取，handler 自行决定是否响应取消
	ctx := goocontext.WithGenerateTraceId(context.Background())
	log := goocontext.Log(ctx).WithTag("goo-kafka-consumer-group", g.id).WithField("msg", m)

	// uniq key
	{
		var uniqKey string
		if key != "" {
			uniqKey = fmt.Sprintf("%s:%s", g.id, key)
		} else {
			uniqKey = fmt.Sprintf("%s:%s:%s", g.id, msg.Topic, goo_utils.MD5([]byte(g.id+msg.Topic+string(msg.Value))))
		}
		if g.cli.redis != nil {
			ok := g.cli.redis.SetNX(uniqKey, goo_utils.M{
				"topic":     msg.Topic,
				"body":      m["body"],
				"headers":   m["headers"],
				"timestamp": m["timestamp"],
			}.String(), 300*time.Second).Val()
			if !ok {
				log.Warn("消息消费失败，并发消费")
				// 去重命中视为已处理，提交 offset，避免卡在重投循环
				session.MarkMessage(msg, "")
				return
			}
			defer func() {
				g.cli.redis.Del(uniqKey)
			}()
		}
	}

	// 建立缓存
	if g.cli.redis != nil && key != "" {
		g.cli.redis.Set(key, goo_utils.M{
			"topic":     msg.Topic,
			"body":      m["body"],
			"headers":   m["headers"],
			"timestamp": m["timestamp"],
		}.String(), time.Hour)
	}

	// 打印日志
	t1 := time.Now()
	defer func() {
		// 删除缓存
		if g.cli.redis != nil && key != "" {
			g.cli.redis.Expire(key, 5*time.Second)
		}

		log = log.WithField("执行时间", fmt.Sprintf("%f", float64(time.Now().Sub(t1).Milliseconds())/1e3))
		if err != nil {
			log.Error("消息消费失败", err)
			return
		}

		log.Debug("消息消费成功")
	}()

	// 执行业务方法
	if err = g.handler(ctx, &ConsumerMessage{ConsumerMessage: msg, GroupSession: session}, nil); err != nil {
		return
	}

	// 提交
	session.MarkMessage(msg, "")

	return
}
