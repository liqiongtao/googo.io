package goo_task_queue

import (
	"time"

	"github.com/go-redis/redis"
)

// requeueOrFail 统一再次入队：达 MaxRetry 进 fail，否则按 nextRunAtMs 写回 pending。
func (t *TaskQueueTasks) requeueOrFail(task *Task, nextRunAtMs int64) error {
	if task.MaxRetry != 0 && task.RetryTimes >= task.MaxRetry {
		return t.taskFail(task)
	}
	return t.requeue(task, nextRunAtMs)
}

func (t *TaskQueueTasks) requeue(task *Task, nextRunAtMs int64) error {
	if nextRunAtMs <= 0 {
		nextRunAtMs = time.Now().UnixMilli()
	}

	score := priorityScore(task.HighPriority, nextRunAtMs)

	pi := t.r.TxPipeline()
	pi.HIncrBy(t.taskInfoKey(task.Id), "retry_times", 1)
	pi.HSet(t.taskInfoKey(task.Id), "ts", nextRunAtMs)
	pi.ZRem(t.TaskProcessingKey, task.Id)
	pi.ZAdd(t.TaskPendingKey, redis.Z{Member: task.Id, Score: score})

	if _, err := pi.Exec(); err != nil {
		t.log().WithTag("requeue").Error(err)
		return err
	}
	return nil
}

func (t *TaskQueueTasks) taskFail(task *Task) error {
	pi := t.r.TxPipeline()
	pi.ZRem(t.TaskProcessingKey, task.Id)
	pi.ZAdd(t.TaskFailKey, redis.Z{Member: task.Id, Score: float64(time.Now().UnixMilli())})

	if _, err := pi.Exec(); err != nil {
		t.log().WithTag("taskFail").Error(err)
		return err
	}
	return nil
}
