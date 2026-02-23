package goo_task_queue

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis"
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
redis.call('ZADD', KEYS[2], ARGV[1], member)

return member

`

	keys := []string{t.TaskPendingKey, t.TaskProcessingKey}
	args := []interface{}{float64(time.Now().Unix())}

	member, err := t.r.Eval(luaScript, keys, args...).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		t.log().WithTag("getOneTask").Error(err)
		return nil, err
	}

	// 任务ID
	taskId := member.(string)

	if taskId == "" || !t.taskExists(taskId) {
		t.r.ZRem(t.TaskPendingKey, member)
		return nil, nil
	}

	task := getTaskByCache(t.r, t.taskInfoKey(taskId))
	if task.Id == "" {
		t.log().WithTag("getOneTask").Warn(fmt.Sprintf("%s not exists", taskId))
		t.taskDel(taskId)
		return nil, nil
	}

	return task, nil
}

func (t *TaskQueueTasks) taskDel(taskIds ...string) error {
	pi := t.r.TxPipeline()

	for _, taskId := range taskIds {
		// 删除任务
		pi.Del(t.taskInfoKey(taskId))
		// 删除待执行队列
		pi.ZRem(t.TaskPendingKey, taskId)
		// 删除执行队列
		pi.ZRem(t.TaskProcessingKey, taskId)
		// 删除失败队列
		pi.ZRem(t.TaskFailKey, taskId)
	}

	if _, err := pi.Exec(); err != nil {
		t.log().WithTag("delTask").Error(err)
		return err
	}

	return nil
}

func (t *TaskQueueTasks) taskExists(taskId string) bool {
	return t.r.Exists(t.taskInfoKey(taskId)).Val() > 0
}

func (t *TaskQueueTasks) log() *goo_log.Entry {
	return goo_log.WithTag("goo-task-queue-tasks", t.pid)
}
