package goo_event

import (
	"sync"

	goo_log "github.com/liqiongtao/googo.io/goo-log"
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
					goo_log.WithTag("goo-event").WithField("topic", topic).Warn("订阅 channel 已满，丢弃消息")
				}
			}
		})
	}
}

// Subscribe 订阅 topic，返回反订阅函数。
// 回调在订阅 goroutine 内同步执行，避免高峰时无界创建 goroutine。
func (ev *Event) Subscribe(topic string, fn SubscribeFunc) (unsubscribe func()) {
	ch := make(chan Message, 64)
	stop := make(chan struct{})

	ev.mu.Lock()
	ev.subscribes[topic] = append(ev.subscribes[topic], ch)
	ev.mu.Unlock()

	goo_utils.AsyncFunc(func() {
		for {
			select {
			case <-goocontext.Root().Done():
				return
			case <-stop:
				return
			case msg := <-ch:
				// 单条回调 panic 不应杀死整个订阅循环
				func() {
					defer func() {
						if r := recover(); r != nil {
							goo_log.WithTag("goo-event").WithField("topic", topic).Error(r)
						}
					}()
					fn(msg)
				}()
			}
		}
	})

	var once sync.Once
	return func() {
		once.Do(func() {
			close(stop)

			ev.mu.Lock()
			defer ev.mu.Unlock()

			chs := ev.subscribes[topic]
			for i, c := range chs {
				if c == ch {
					ev.subscribes[topic] = append(chs[:i], chs[i+1:]...)
					break
				}
			}
			if len(ev.subscribes[topic]) == 0 {
				delete(ev.subscribes, topic)
			}
		})
	}
}
