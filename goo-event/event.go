package goo_event

import (
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
	"sync"
)

type Event struct {
	subscribes map[string][]MessageChan
	mu         sync.RWMutex
}

func New() *Event {
	return &Event{subscribes: map[string][]MessageChan{}}
}

// 发布
func (ev *Event) Publish(topic string, data interface{}) {
	ev.mu.RLock()
	defer ev.mu.RUnlock()

	if chs, ok := ev.subscribes[topic]; ok {
		channels := append([]MessageChan{}, chs...)
		goo_utils.AsyncFunc(func() {
			for _, ch := range channels {
				ch <- Message{Topic: topic, Data: data}
			}
		})
	}
}

// 订阅
func (ev *Event) Subscribe(topic string, fn SubscribeFunc) {
	ev.mu.Lock()
	defer ev.mu.Unlock()

	if _, ok := ev.subscribes[topic]; !ok {
		ev.subscribes[topic] = []MessageChan{}
	}

	ch := make(chan Message)
	ev.subscribes[topic] = append(ev.subscribes[topic], ch)

	goo_utils.AsyncFunc(func() {
		for {
			select {
			case msg := <-ch:
				fn(msg)
			}
		}
	})
}
