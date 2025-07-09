package goo_task_queue

import (
	"github.com/go-redis/redis"
	goo_context "github.com/liqiongtao/googo.io/goo-context"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
	"time"
)

type TaskQueueLeader struct {
	*TaskQueue
}

func (l *TaskQueueLeader) Generate() {
	for {
		if l.handler() {
			return
		}
		time.Sleep(time.Second)
	}
}

func (l *TaskQueueLeader) handler() bool {
	if !l.lock() {
		return false
	}
	defer l.unlock()

	ctx := goo_context.WithCancel()

	goo_utils.AsyncFuncGroup(func() {
		for {
			select {
			case <-ctx.Done():
				return

			default:
				l.expire()
				time.Sleep(time.Second)
			}
		}
	}, func() {
		for {
			select {
			case <-ctx.Done():
				return

			default:
				l.recover()
				time.Sleep(time.Second)
			}
		}
	})

	return true
}

// 恢复
func (l *TaskQueueLeader) recover() error {
	zs := l.r.ZRangeWithScores(l.TaskProcessingKey, 0, -1).Val()
	if len(zs) == 0 {
		time.Sleep(time.Second * 10) // 没有任务时休眠
		return nil
	}

	for _, z := range zs {
		taskId := z.Member.(string)

		if taskId == "" || !l.taskExists(taskId) {
			l.r.ZRem(l.TaskProcessingKey, z.Member)
			continue
		}

		task := getTaskByCache(l.r, l.taskInfoKey(taskId))
		if task.Id == "" {
			l.TaskQueueTasks.taskDel(taskId)
			continue
		}

		if ts := time.Now().Unix() - int64(z.Score); ts > task.Timeout {
			pi := l.r.TxPipeline()

			// 删除执行队列
			pi.ZRem(l.TaskProcessingKey, taskId)
			// 添加待执行队列
			pi.ZAdd(l.TaskPendingKey, redis.Z{Member: taskId, Score: float64(time.Now().Unix() + 60)})

			if _, err := pi.Exec(); err != nil {
				l.log().WithTag("retry").Error(err)
			}
		}
	}

	return nil
}

func (l *TaskQueueLeader) lock() bool {
	return l.r.SetNX(l.TaskLeadLockKey, time.Now().Unix(), time.Second*10).Val()
}

func (l *TaskQueueLeader) unlock() error {
	return l.r.Del(l.TaskLeadLockKey).Err()
}

func (l *TaskQueueLeader) expire() bool {
	return l.r.Expire(l.TaskLeadLockKey, time.Second*10).Val()
}

func (l *TaskQueueLeader) log() *goo_log.Entry {
	return goo_log.WithTag("goo-task-queue-leader")
}
