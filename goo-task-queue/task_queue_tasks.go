package goo_task_queue

import (
	"fmt"
	"github.com/go-redis/redis"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"time"
)

type TaskQueueTasks struct {
	*TaskQueue
}

func (t *TaskQueueTasks) getOneTask() (*Task, error) {
	// 加锁
	if !t.lock() {
		return nil, nil
	}
	defer t.unlock()

	// 获取待执行队列
	zs := t.r.ZRangeWithScores(t.TaskPendingKey, 0, time.Now().Unix()).Val()
	if len(zs) == 0 {
		return nil, nil
	}

	// 任务ID
	taskId := zs[0].Member.(string)

	// 任务不存在
	if taskId == "" || !t.taskExists(taskId) {
		t.log().WithTag("getOneTask").Warn(fmt.Sprintf("%s not exists", taskId))
		t.taskDel(taskId)
		return nil, nil
	}

	task := getTaskByCache(t.r, t.taskInfoKey(taskId))

	// 执行队列
	{
		pi := t.r.TxPipeline()

		// 删除待执行队列
		pi.ZRem(t.TaskPendingKey, zs[0].Member)
		// 添加执行队列
		pi.ZAdd(t.TaskProcessingKey, redis.Z{Member: zs[0].Member, Score: float64(time.Now().Unix())})

		if _, err := pi.Exec(); err != nil {
			t.log().WithTag("getOneTask").Error("add taskProcessingKey fail", err)
			return nil, err
		}
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

func (t *TaskQueueTasks) lock() bool {
	return t.r.SetNX(t.TaskGetLockKey, time.Now().Unix(), time.Second*10).Val()
}

func (t *TaskQueueTasks) unlock() error {
	return t.r.Del(t.TaskGetLockKey).Err()
}

func (t *TaskQueueTasks) log() *goo_log.Entry {
	return t.TaskQueue.log().WithTag("goo-task-queue-tasks")
}
