package goo_log

type KafkaAdapter struct {
	topic string
}

func NewKafkaAdapter(topic string) *KafkaAdapter {
	return &KafkaAdapter{
		topic: topic,
	}
}

func (ca KafkaAdapter) Write(msg *Message) {
	// todo:: 1. 包冲突问题； 2. 初始化问题

	//_, _, err := goo_kafka.Producer().SendAsyncMessage(ca.topic, ca.topic, msg.JSON())
	//if err != nil {
	//	log.Println(err.Error())
	//}
}
