package goo_task_queue

import (
	"encoding/json"
	"time"

	goo_redis "github.com/liqiongtao/googo.io/goo-redis"
)

const (
	defaultMaxRetry    = 99
	defaultTaskTimeout = int64(2 * time.Hour / time.Millisecond) // 毫秒，默认 2 小时
)

type Task struct {
	Id           string `json:"id,omitempty"`            // 任务唯一ID 必填
	Type         string `json:"type,omitempty"`          // 任务类型 必填
	Payload      string `json:"payload,omitempty"`       // 任务数据 必填
	HighPriority int    `json:"high_priority,omitempty"` // 任务优先级 非必填 1=高优先级, 0=低优先级
	MaxRetry     int    `json:"max_retry,omitempty"`     // 最大重试次数 非必填 默认99次
	RetryTimes   int    `json:"retry_times,omitempty"`   // 重试次数 非必填 默认0次
	Timeout      int64  `json:"timeout,omitempty"`       // 任务超时时间（毫秒）非必填 默认2小时
	Ts           int64  `json:"ts,omitempty"`            // 调度时间（毫秒时间戳），可设为未来表示延迟执行
	Generation   int64  `json:"generation,omitempty"`    // 执行代数（抢占时分配，收尾校验用，业务勿写）
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
	task.Ts, _ = r.HGet(key, "ts").Int64()
	task.Generation, _ = r.HGet(key, "generation").Int64()

	if task.Ts == 0 {
		task.Ts = time.Now().UnixMilli()
	}

	if task.MaxRetry == 0 {
		task.MaxRetry = defaultMaxRetry
	}
	if task.Timeout == 0 {
		task.Timeout = defaultTaskTimeout
	}

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
		"ts":            t.Ts,
	}
}

func (t *Task) Json() []byte {
	b, _ := json.Marshal(&t)
	return b
}

func (t *Task) String() string {
	return string(t.Json())
}
