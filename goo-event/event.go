package goo_event

import (
	"sync"

	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
	"github.com/liqiongtao/googo.io/goocontext"
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
		msg := Message{Topic: topic, Data: data}
		goo_utils.AsyncFunc(func() {
			for _, ch := range channels {
				select {
				case ch <- msg:
				default:
					// 订阅方处理过慢时丢弃，避免永久阻塞发布 goroutine
				}
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

	ch := make(chan Message, 64)
	ev.subscribes[topic] = append(ev.subscribes[topic], ch)

	goo_utils.AsyncFunc(func() {
		for {
			select {
			case <-goocontext.Root().Done():
				return
			case msg := <-ch:
				goo_utils.AsyncFunc(func() {
					fn(msg)
				})
			}
		}
	})
}
