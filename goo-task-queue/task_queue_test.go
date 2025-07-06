package goo_task_queue

import (
	"context"
	"errors"
	"fmt"
	goo_redis "github.com/liqiongtao/googo.io/goo-redis"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
	"testing"
	"time"
)

var (
	redisConfig = goo_redis.Config{
		Addr:   "127.0.0.1:6379",
		DB:     0,
		Prefix: "test",
	}
)

func TestTaskQueue_Publish(t *testing.T) {
	r, err := goo_redis.New(redisConfig)
	if err != nil {
		fmt.Println(err)
		return
	}

	q := New(r)

	for i := range 3 {
		go func() {
			task := &Task{
				Id:      fmt.Sprintf("%d", i),
				Type:    "test",
				Payload: goo_utils.M{"name": fmt.Sprintf("name-%d", i)}.String(),
				//MaxRetry:     3,
				//HighPriority: 1,
			}
			if i == 2 {
				task.HighPriority = 1
			}
			q.Publish(task)
		}()
		if i%3 == 0 {
			time.Sleep(time.Second * 1)
		}
	}

	time.Sleep(time.Second * 1)
}

func TestTaskQueue_Subscribe(t *testing.T) {
	r, err := goo_redis.New(redisConfig)
	if err != nil {
		fmt.Println(err)
		return
	}

	q := New(r)

	q.Subscribe(1, func(ctx context.Context, value any) error {
		time.Sleep(time.Second * 1)
		return errors.New("test error")
	})
}

func TestTaskQueue_Subscribe2(t *testing.T) {
	r, err := goo_redis.New(redisConfig)
	if err != nil {
		fmt.Println(err)
		return
	}

	q := New(r)

	q.Subscribe(1, func(ctx context.Context, value any) error {
		time.Sleep(time.Second * 1)
		return nil
	})
}
