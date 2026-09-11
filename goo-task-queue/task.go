package gootaskqueue

import (
	"encoding/json"
	"strconv"
	"time"

	gooredis "github.com/liqiongtao/googo.io/goo-redis"
)

const (
	defaultMaxRetry    = 99
	defaultTaskTimeout = int64(2 * time.Hour / time.Millisecond) // 毫秒，默认 2 小时
	taskInfoTTLSec     = int64(48 * 60 * 60)
)

// infoTTLSecUntil：基础 48h + 距 nextRun 的延迟，避免长延迟任务 info 先过期被清掉
func infoTTLSecUntil(nextRunAtMs int64) int64 {
	ttl := taskInfoTTLSec
	if delaySec := (nextRunAtMs - time.Now().UnixMilli()) / 1000; delaySec > 0 {
		ttl += delaySec
	}
	return ttl
}

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

func getTaskByCache(r *gooredis.Client, key string) (*Task, error) {
	m, err := r.HGetAll(key).Result()
	if err != nil {
		return nil, err
	}
	if len(m) == 0 {
		return &Task{}, nil
	}

	task := &Task{
		Id:      m["id"],
		Type:    m["type"],
		Payload: m["payload"],
	}

	task.HighPriority, _ = strconv.Atoi(m["high_priority"])
	task.MaxRetry, _ = strconv.Atoi(m["max_retry"])
	task.RetryTimes, _ = strconv.Atoi(m["retry_times"])
	task.Timeout, _ = strconv.ParseInt(m["timeout"], 10, 64)
	task.Ts, _ = strconv.ParseInt(m["ts"], 10, 64)
	task.Generation, _ = strconv.ParseInt(m["generation"], 10, 64)

	if task.Ts == 0 {
		task.Ts = time.Now().UnixMilli()
	}

	if task.MaxRetry == 0 {
		task.MaxRetry = defaultMaxRetry
	}
	if task.Timeout == 0 {
		task.Timeout = defaultTaskTimeout
	}

	return task, nil
}

func (t *Task) MapData() map[string]any {
	return map[string]any{
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
