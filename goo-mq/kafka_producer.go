package gooMQ

import (
	"fmt"
	"github.com/Shopify/sarama"
	gooLog "github.com/liqiongtao/googo.io/goo/log"
	"time"
)

type KafkaProducer struct {
	*Kafka
	producer sarama.AsyncProducer
}

func (*KafkaProducer) config() *sarama.Config {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.Timeout = 5 * time.Second
	config.Version = sarama.V0_10_2_0
	return config
}

func (p *KafkaProducer) Init() {
	gooLog.Debug("[kafka-producer-init]")

	producer, err := sarama.NewAsyncProducer(p.Addrs, p.config())
	if err != nil {
		gooLog.Debug("[kafka-producer-error]", err.Error())
		panic(err.Error())
	}

	go func() {
		for {
			select {
			case suc := <-producer.Successes():
				gooLog.Debug("[kafka-producer-success]",
					fmt.Sprintf("partitions=%d topic=%s offset=%d value=%s",
						suc.Partition, suc.Topic, suc.Offset, suc.Value))

			case err := <-producer.Errors():
				gooLog.Error("[kafka-producer-error]", err.Error())

			case <-p.Context.Done():
				return
			}
		}
	}()

	p.producer = producer
}

func (p *KafkaProducer) SendMessage(topic string, message []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(message),
		Key:   sarama.StringEncoder(fmt.Sprintf("%d", time.Now().UnixNano())),
	}

	p.producer.Input() <- msg

	return nil
}
