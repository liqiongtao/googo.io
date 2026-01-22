package goo_task_queue

import (
	"math/rand"
	"time"

	"github.com/go-redis/redis"
	goo_context "github.com/liqiongtao/googo.io/goo-context"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
)

type TaskQueueLeader struct {
	*TaskQueue
}

func (l *TaskQueueLeader) Generate() {
	for {
		if l.handler() {
			return
		}
		time.Sleep(time.Duration(rand.Intn(600)+200) * time.Millisecond)
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
	}, func() {
		for {
			select {
			case <-ctx.Done():
				return

			default:
				l.workers()
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
			pi.ZAdd(l.TaskPendingKey, redis.Z{Member: taskId, Score: float64(task.Ts)})

			if _, err := pi.Exec(); err != nil {
				l.log().WithTag("retry").Error(err)
			}
		}
	}

	return nil
}

// 处理超时的workers
func (l *TaskQueueLeader) workers() error {
	workerIds := l.r.HKeys(l.TaskWorkersKey).Val()
	if len(workerIds) == 0 {
		time.Sleep(time.Second * 10) // 没有任务时休眠
		return nil
	}

	for _, workerId := range workerIds {
		str := l.r.HGet(l.TaskWorkersKey, workerId).Val()
		if str == "" {
			continue
		}
		ti, err := time.ParseInLocation("2006-01-02 15:04:05", str, time.Local)
		if err != nil {
			continue
		}
		if time.Now().Unix()-ti.Unix() > 20 {
			l.r.HDel(l.TaskWorkersKey, workerId)
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
	return goo_log.WithTag("goo-task-queue-leader", l.pid)
}
