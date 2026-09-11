package goo_event

import "sync"

var (
	__event *Event
	__once  sync.Once
)

func ensureEvent() *Event {
	__once.Do(func() {
		__event = New()
	})
	return __event
}

func Default() *Event {
	return ensureEvent()
}

func Publish(topic string, data any) {
	ensureEvent().Publish(topic, data)
}

func Subscribe(topic string, fn SubscribeFunc) (unsubscribe func()) {
	return ensureEvent().Subscribe(topic, fn)
}
