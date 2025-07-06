package goo_task_queue

import (
	"context"
	"encoding/json"
	"time"
)

type TaskQueueHandler func(ctx context.Context, value any) error

type Task struct {
	Id                string `json:"id,omitempty"`                  // 任务唯一ID 必填
	Type              string `json:"type,omitempty"`                // 任务类型 必填
	Payload           string `json:"payload,omitempty"`             // 任务数据 必填
	HighPriority      int    `json:"high_priority,omitempty"`       // 任务优先级 非必填 1=高优先级, 0=低优先级
	MaxRetry          int    `json:"max_retry,omitempty"`           // 最大重试次数 非必填 默认99次
	RetryTimes        int    `json:"retry_times,omitempty"`         // 重试次数 非必填 默认0次
	RetryIntervalTime int    `json:"retry_interval_time,omitempty"` // 重试间隔时间 非必填 默认1秒
	Timeout           int64  `json:"timeout,omitempty"`             // 任务超时时间 非必填 默认30分钟
}

func (t *Task) MapData() map[string]interface{} {
	return map[string]interface{}{
		"id":                  t.Id,
		"type":                t.Type,
		"payload":             t.Payload,
		"high_priority":       t.HighPriority,
		"max_retry":           t.MaxRetry,
		"retry_times":         t.RetryTimes,
		"retry_interval_time": t.RetryIntervalTime,
		"timeout":             t.Timeout,
	}
}

func (t *Task) Score() float64 {
	if t.HighPriority == 1 {
		return 1
	}
	return float64(time.Now().Unix()) + float64(t.RetryTimes*60)
}

func (t *Task) Json() []byte {
	b, _ := json.Marshal(&t)
	return b
}

func (t *Task) String() string {
	return string(t.Json())
}
