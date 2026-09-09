package goo_task_queue

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/liqiongtao/googo.io/goo-redis"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
)

type TaskQueueTasks struct {
	*TaskQueue
}

func (t *TaskQueueTasks) getOneTask() (*Task, error) {
	luaScript := `
local tasks = redis.call('ZRANGEBYSCORE', KEYS[1], 0, ARGV[1], 'LIMIT', 0, 1)
if #tasks == 0 then
	return nil
end

local member = tasks[1]

redis.call('ZREM', KEYS[1], member)
redis.call('ZADD', KEYS[2], ARGV[2], member)

local gen = redis.call('HINCRBY', ARGV[3] .. member, 'generation', 1)
return {member, tostring(gen)}
`

	nowMs := float64(time.Now().UnixMilli())
	keys := []string{t.TaskPendingKey, t.TaskProcessingKey}
	// ARGV[1]=可调度上限；ARGV[2]=processing 开始时间；ARGV[3]=info key 前缀
	args := []any{nowMs + 0.5, nowMs, t.TaskInfoKey + ":"}

	result, err := t.r.Eval(luaScript, keys, args...).Result()
	if err != nil {
		if errors.Is(err, goo_redis.ErrNil) {
			return nil, nil
		}
		t.log().WithTag("getOneTask").Error(err)
		return nil, err
	}

	arr, ok := result.([]any)
	if !ok || len(arr) < 2 {
		t.log().WithTag("getOneTask").ErrorF("unexpected lua result: %v", result)
		return nil, fmt.Errorf("unexpected getOneTask result")
	}

	taskId, _ := arr[0].(string)
	genStr, _ := arr[1].(string)

	if taskId == "" || !t.taskExists(taskId) {
		t.r.ZRem(t.TaskPendingKey, taskId)
		t.r.ZRem(t.TaskProcessingKey, taskId)
		return nil, nil
	}

	task := getTaskByCache(t.r, t.taskInfoKey(taskId))
	if task.Id == "" {
		t.log().WithTag("getOneTask").Warn(fmt.Sprintf("%s not exists", taskId))
		_ = t.taskDelForce(taskId)
		return nil, nil
	}

	gen, err := strconv.ParseInt(genStr, 10, 64)
	if err != nil {
		t.log().WithTag("getOneTask").Error(err)
		return nil, err
	}
	task.Generation = gen
	return task, nil
}

func (t *TaskQueueTasks) taskExists(taskId string) bool {
	return t.r.Exists(t.taskInfoKey(taskId)).Val() > 0
}

func (t *TaskQueueTasks) log() *goo_log.Entry {
	return goo_log.WithTag("goo-task-queue-tasks", t.pid)
}
