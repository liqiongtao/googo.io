package goo_event

type Message struct {
	Topic string
	Data  any
}

type MessageChan chan Message

type SubscribeFunc func(msg Message)
