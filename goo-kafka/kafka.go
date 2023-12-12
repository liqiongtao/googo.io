package goo_kafka

import (
	"github.com/IBM/sarama"
	goo_context "github.com/liqiongtao/googo.io/goo-context"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
)

var (
	__client *client
)

// 初始化
func Init(conf Config) error {
	__client = &client{conf: conf}
	goo_utils.AsyncFunc(func() {
		select {
		case <-goo_context.Cancel().Done():
			__client.Close()
			return
		}
	})
	return __client.init()
}

// 客户端
func Client() *client {
	return __client
}

// 生产者
func Producer() iProducer {
	return &producer{client: __client, msg: &sarama.ProducerMessage{}}
}

// 消费者
func Consumer() iConsumer {
	return &consumer{client: __client}
}

// 主题列表
func Topics() []string {
	if __client == nil {
		return []string{}
	}

	topics, err := __client.Topics()
	if err != nil {
		goo_log.WithTag("goo-kafka").Error(err)
		return []string{}
	}

	return topics
}

// 分区数量
func Partitions(topic string) []int32 {
	if __client == nil {
		return []int32{}
	}

	partitions, err := __client.Partitions(topic)
	if err != nil {
		goo_log.WithTag("goo-kafka").WithField("topic", topic).Error(err)
		return []int32{}
	}

	return partitions
}

// 分区数量
func OffsetInfo(topic, groupId string) (data []map[string]int64) {
	data = []map[string]int64{}

	if __client == nil {
		return
	}

	partitions := Partitions(topic)
	if l := len(partitions); l == 0 {
		return
	}

	var (
		l = goo_log.WithTag("goo-kafka").WithField("groupId", groupId).WithField("topic", topic)
	)

	om, err := sarama.NewOffsetManagerFromClient(groupId, __client.Client)
	if err != nil {
		l.Error(err)
		return
	}

	for _, partition := range partitions {
		offset, err := __client.GetOffset(topic, partition, -1)
		if err != nil {
			l.Error(err)
			continue
		}

		pom, err := om.ManagePartition(topic, partition)
		if err != nil {
			l.Error(err)
			continue
		}

		nextOffset, msg := pom.NextOffset()
		if msg != "" {
			l.Error(msg)
			continue
		}

		backlog := offset
		if nextOffset != -1 {
			backlog -= nextOffset
		}

		data = append(data, map[string]int64{
			"partition":  int64(partition),
			"offset":     offset,
			"nextOffset": nextOffset,
			"backlog":    backlog,
		})
	}

	return
}
