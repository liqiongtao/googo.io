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

	taskId := luaResultString(arr[0])
	genStr := luaResultString(arr[1])
	if taskId == "" {
		return nil, nil
	}

	exists, err := t.taskExists(taskId)
	if err != nil {
		t.log().WithTag("getOneTask").Error(err)
		t.putBackOnReadError(taskId, genStr)
		return nil, err
	}
	if !exists {
		t.r.ZRem(t.TaskPendingKey, taskId)
		t.r.ZRem(t.TaskProcessingKey, taskId)
		return nil, nil
	}

	task, err := getTaskByCache(t.r, t.taskInfoKey(taskId))
	if err != nil {
		t.log().WithTag("getOneTask").Error(err)
		t.putBackOnReadError(taskId, genStr)
		return nil, err
	}
	if task.Id == "" {
		t.log().WithTag("getOneTask").Warn(fmt.Sprintf("%s not exists", taskId))
		_ = t.taskDelForce(taskId)
		return nil, nil
	}

	gen, err := strconv.ParseInt(genStr, 10, 64)
	if err != nil {
		t.log().WithTag("getOneTask").Error(err)
		// Lua 已将任务移入 processing，解析失败时放回 pending，避免卡住至超时回收
		if pbErr := t.putBack(task); pbErr != nil {
			t.log().WithTag("getOneTask").WithField("task_id", taskId).ErrorF("putBack after parse error: %v", pbErr)
		}
		return nil, err
	}
	task.Generation = gen
	return task, nil
}

func (t *TaskQueueTasks) putBackOnReadError(taskId, genStr string) {
	gen, _ := strconv.ParseInt(genStr, 10, 64)
	task := &Task{Id: taskId, Generation: gen}
	if pbErr := t.putBack(task); pbErr != nil {
		t.log().WithTag("getOneTask").WithField("task_id", taskId).ErrorF("putBack after redis read error: %v", pbErr)
	}
}

func luaResultString(v any) string {
	switch x := v.(type) {
	case string:
		return x
	case []byte:
		return string(x)
	case nil:
		return ""
	default:
		return fmt.Sprint(x)
	}
}

func (t *TaskQueueTasks) taskExists(taskId string) (bool, error) {
	n, err := t.r.Exists(t.taskInfoKey(taskId)).Result()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

func (t *TaskQueueTasks) log() *goo_log.Entry {
	return goo_log.WithTag("goo-task-queue-tasks", t.pid)
}
