package goo_task_queue

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	goo_log "github.com/liqiongtao/googo.io/goo-log"
	goo_redis "github.com/liqiongtao/googo.io/goo-redis"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
)

type TaskQueue struct {
	pid        string
	instanceId string
	r          *goo_redis.Client

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

	localIp, err := goo_utils.LocalIP()
	if err != nil || localIp == "" {
		if host, herr := os.Hostname(); herr == nil && host != "" {
			localIp = host
		} else {
			localIp = "unknown"
		}
	}
	pid := os.Getpid()
	q := &TaskQueue{
		pid:              strconv.Itoa(pid),
		instanceId:       fmt.Sprintf("%s:%d", localIp, pid),
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
	// <=0 表示关闭内存门控；>0 为占用上限百分比
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
	return goo_log.WithTag("goo-task-queue", q.pid)
}
