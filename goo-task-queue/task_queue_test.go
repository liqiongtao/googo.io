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

	for i := range 10000 {
		task := &Task{
			Id:      fmt.Sprintf("%d", i+1),
			Type:    "test",
			Payload: goo_utils.M{"name": fmt.Sprintf("name-%d", i+1)}.String(),
		}
		q.Publish(task)
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

	q.Subscribe(3, func(ctx context.Context, task *Task) error {
		time.Sleep(time.Second * 2)
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

	q.Subscribe(3, func(ctx context.Context, task *Task) error {
		time.Sleep(time.Second)
		return nil
	})
}

func TestTaskQueue_TaskCount(t *testing.T) {
	r, err := goo_redis.New(redisConfig)
	if err != nil {
		fmt.Println(err)
		return
	}

	q := New(r)

	fmt.Println("PendingCount:", q.PendingCount())
	fmt.Println("ProcessingCount:", q.ProcessingCount())
	fmt.Println("FailCount:", q.FailCount())
}
