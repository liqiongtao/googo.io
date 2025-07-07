package goo_task_queue

import (
	"encoding/json"
	goo_redis "github.com/liqiongtao/googo.io/goo-redis"
	"time"
)

type Task struct {
	Id           string `json:"id,omitempty"`            // 任务唯一ID 必填
	Type         string `json:"type,omitempty"`          // 任务类型 必填
	Payload      string `json:"payload,omitempty"`       // 任务数据 必填
	HighPriority int    `json:"high_priority,omitempty"` // 任务优先级 非必填 1=高优先级, 0=低优先级
	MaxRetry     int    `json:"max_retry,omitempty"`     // 最大重试次数 非必填 默认99次
	RetryTimes   int    `json:"retry_times,omitempty"`   // 重试次数 非必填 默认0次
	Timeout      int64  `json:"timeout,omitempty"`       // 任务超时时间 非必填 默认30分钟
}

func getTaskByCache(r *goo_redis.Client, key string) *Task {
	task := &Task{
		Id:      r.HGet(key, "id").Val(),
		Type:    r.HGet(key, "type").Val(),
		Payload: r.HGet(key, "payload").Val(),
	}

	task.HighPriority, _ = r.HGet(key, "high_priority").Int()
	task.MaxRetry, _ = r.HGet(key, "max_retry").Int()
	task.RetryTimes, _ = r.HGet(key, "retry_times").Int()
	task.Timeout, _ = r.HGet(key, "timeout").Int64()

	return task
}

func (t *Task) MapData() map[string]interface{} {
	return map[string]interface{}{
		"id":            t.Id,
		"type":          t.Type,
		"payload":       t.Payload,
		"high_priority": t.HighPriority,
		"max_retry":     t.MaxRetry,
		"retry_times":   t.RetryTimes,
		"timeout":       t.Timeout,
	}
}

func (t *Task) Score() float64 {
	if t.HighPriority == 1 {
		return 1
	}
	return float64(time.Now().Unix()) + 60
}

func (t *Task) Json() []byte {
	b, _ := json.Marshal(&t)
	return b
}

func (t *Task) String() string {
	return string(t.Json())
}
