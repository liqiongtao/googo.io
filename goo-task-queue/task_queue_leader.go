package goo_task_queue

import (
	"context"
	"math/rand"
	"time"

	goo_log "github.com/liqiongtao/googo.io/goo-log"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
)

type TaskQueueLeader struct {
	*TaskQueue
}

func (l *TaskQueueLeader) Generate(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		if l.handler(ctx) {
			return
		}

		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Duration(rand.Intn(600)+200) * time.Millisecond):
		}
	}
}

func (l *TaskQueueLeader) handler(ctx context.Context) bool {
	if !l.lock() {
		return false
	}
	defer l.unlock()

	goo_utils.AsyncFuncGroup(func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				l.expire()
				if !sleepOrDone(ctx, time.Second) {
					return
				}
			}
		}
	}, func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				l.recover(ctx)
				if !sleepOrDone(ctx, time.Second) {
					return
				}
			}
		}
	}, func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				l.workers(ctx)
				if !sleepOrDone(ctx, time.Second) {
					return
				}
			}
		}
	})

	return true
}

// recover 按名次从最老开始分批扫描 processing，避免同 score exclusive 游标跳过成员
func (l *TaskQueueLeader) recover(ctx context.Context) error {
	const batch int64 = 100
	var start int64
	scanned := 0
	nowMs := time.Now().UnixMilli()

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		nowMs = time.Now().UnixMilli()
		zs := l.r.ZRangeWithScores(l.TaskProcessingKey, start, start+batch-1).Val()
		if len(zs) == 0 {
			break
		}
		scanned += len(zs)

		requeued := 0
		for _, z := range zs {
			taskId, _ := z.Member.(string)
			startMs := int64(z.Score)

			if taskId == "" || !l.taskExists(taskId) {
				l.r.ZRem(l.TaskProcessingKey, z.Member)
				requeued++
				continue
			}

			task := getTaskByCache(l.r, l.taskInfoKey(taskId))
			if task.Id == "" {
				_ = l.TaskQueueTasks.taskDelForce(taskId)
				requeued++
				continue
			}

			if nowMs-startMs > task.Timeout {
				if err := l.TaskQueueTasks.requeueOrFail(task, time.Now().UnixMilli()); err != nil {
					l.log().WithTag("recover").Error(err)
					continue
				}
				requeued++
			}
		}

		if requeued > 0 {
			// 有删除，名次前移，从当前 start 重扫
			continue
		}
		if int64(len(zs)) < batch {
			break
		}
		// 本批均未超时：名次后移（同 score 成员连续排列，不会被跳过）
		start += batch
	}

	if scanned == 0 {
		sleepOrDone(ctx, time.Second*10)
	}
	return nil
}

func (l *TaskQueueLeader) workers(ctx context.Context) error {
	workerIds := l.r.HKeys(l.TaskWorkersKey).Val()
	if len(workerIds) == 0 {
		sleepOrDone(ctx, time.Second*10)
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
