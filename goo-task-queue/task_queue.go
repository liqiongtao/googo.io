package goo_task_queue

import (
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	goo_redis "github.com/liqiongtao/googo.io/goo-redis"
	"math/rand"
	"time"
)

type TaskQueue struct {
	r *goo_redis.Client

	*TaskQueueKeys
	*TaskQueueCount
	*TaskQueueList
	*TaskQueueTasks
	*TaskQueueLeader
	*TaskQueuePublisher
	*TaskQueueSubscriber

	// 最大内存占用百分比
	MaxMemoryPercent float64
}

func New(r *goo_redis.Client) *TaskQueue {
	rand.Seed(time.Now().UnixNano())

	keys := NewTaskQueueKeys()
	if r.Config.Prefix != "" {
		keys.WithPrefix(r.Config.Prefix)
	}

	q := &TaskQueue{
		r:                r,
		TaskQueueKeys:    keys,
		MaxMemoryPercent: 90,
	}

	q.TaskQueueCount = &TaskQueueCount{TaskQueue: q}
	q.TaskQueueList = &TaskQueueList{TaskQueue: q}
	q.TaskQueueTasks = &TaskQueueTasks{TaskQueue: q}
	q.TaskQueueLeader = &TaskQueueLeader{TaskQueue: q}
	q.TaskQueuePublisher = &TaskQueuePublisher{TaskQueue: q}
	q.TaskQueueSubscriber = &TaskQueueSubscriber{TaskQueue: q}

	return q
}

func (q *TaskQueue) WithMaxMemoryPercent(percent float64) *TaskQueue {
	q.MaxMemoryPercent = percent
	return q
}

func (q *TaskQueue) Publish(tasks ...*Task) error {
	return q.TaskQueuePublisher.Publish(tasks...)
}

func (q *TaskQueue) Subscribe(limit int, handler TaskQueueHandler) {
	q.TaskQueueSubscriber.Subscribe(limit, handler)
}

func (q *TaskQueue) log() *goo_log.Entry {
	return goo_log.WithTag("goo-task-queue")
}
