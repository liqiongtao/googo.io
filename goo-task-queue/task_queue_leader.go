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
	if !l.lock() {
		return
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

		if taskId == "" || !l.TaskQueueTasks.taskExists(taskId) {
			l.TaskQueueTasks.taskDel(taskId)
			continue
		}

		task := getTaskByCache(l.r, l.taskInfoKey(taskId))

		if ts := time.Now().Unix() - int64(z.Score); ts > task.Timeout {
			pi := l.r.TxPipeline()

			// 删除执行队列
			pi.ZRem(l.TaskProcessingKey, task.Id)
			// 添加待执行队列
			pi.ZAdd(l.TaskPendingKey, redis.Z{Member: task.Id, Score: float64(time.Now().Unix() + 60)})

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
	return l.TaskQueue.log().WithTag("goo-task-queue-leader")
}
