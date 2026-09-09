package goo_task_queue

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/liqiongtao/googo.io/goo"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
	"github.com/liqiongtao/googo.io/goocontext"
)

type TaskQueueHandler func(ctx context.Context, task *Task) error

type TaskQueueSubscriber struct {
	*TaskQueue
}

func (s *TaskQueueSubscriber) Subscribe(limit int, handler TaskQueueHandler) {
	if s.r == nil {
		err := errors.New("redis未初始化")
		s.log().Error(err)
		return
	}

	// 进程级退出 context：Leader / 心跳 / 拉取共用，避免多次 WithCancel
	ctx := goocontext.Root()

	goo_utils.AsyncFunc(func() {
		s.TaskQueueLeader.Generate(ctx)
	})

	goo_utils.AsyncFunc(func() {
		s.heartBeat(ctx)
	})

	if limit <= 0 {
		limit = runtime.NumCPU() * 2
	}

	s.log().InfoF("任务监听成功 并发数=%d workerId=%s", limit, s.workId())

	var (
		limitCH = make(chan any, limit)
		taskCH  = make(chan *Task, limit)
		done    = make(chan any)
	)

	// 执行任务
	goo_utils.AsyncFunc(func() {
		defer func() { done <- struct{}{} }()

		var (
			wg  sync.WaitGroup
			cnt int64
		)

		for {
			task, ok := <-taskCH
			if !ok {
				break
			}

			atomic.AddInt64(&cnt, 1)
			wg.Add(1)

			goo_utils.AsyncFunc(func() {
				defer func() {
					<-limitCH
					left := atomic.AddInt64(&cnt, -1)
					wg.Done()
					s.log().InfoF("剩余 %d 任务正在执行", left)
				}()

				s.taskHandle(task, handler)
			})
		}

		wg.Wait()
		s.log().Info("全部任务执行完毕，执行退出")
	})

	// 获取任务：同步拉取，空队列时不占用执行槽睡眠
	goo_utils.AsyncFunc(func() {
		defer func() { done <- struct{}{} }()

		for {
			select {
			case <-ctx.Done():
				s.log().Info("准备退出")
				close(taskCH)
				return
			default:
			}

			percent, err := goo_utils.MemoryUsedPercent()
			if err != nil || percent >= s.MaxMemoryPercent {
				n := rand.Intn(600) + 3000
				s.log().WarnF("内存占用超过最大限制=%0.2f%%，等待%dms后重试", percent, n)
				if !sleepOrDone(ctx, time.Duration(n)*time.Millisecond) {
					s.log().Info("准备退出")
					close(taskCH)
					return
				}
				continue
			}

			// 先占执行槽再拉取，保证本 Worker 在途任务数不超过 limit
			select {
			case <-ctx.Done():
				s.log().Info("准备退出")
				close(taskCH)
				return
			case limitCH <- struct{}{}:
			}

			task, err := s.getOneTask()
			if err != nil || task == nil {
				<-limitCH // 立即归还，空载睡眠不占槽
				if !sleepOrDone(ctx, time.Duration(rand.Intn(600)+200)*time.Millisecond) {
					s.log().Info("准备退出")
					close(taskCH)
					return
				}
				continue
			}

			select {
			case <-ctx.Done():
				<-limitCH
				// 已抢占但尚未交给执行侧：放回 pending，避免卡在 Timeout 才被 Leader 回收
				_ = s.TaskQueueTasks.putBack(task)
				s.log().Info("准备退出")
				close(taskCH)
				return
			case taskCH <- task:
			}
		}
	})

	for i := 0; i < 2; i++ {
		<-done
	}

	goo_utils.AsyncFuncGroup(func() {
		s.r.HDel(s.TaskWorkersKey, s.workId())
	}, func() {
		close(done)
	})
}

// sleepOrDone 睡眠，若 ctx 取消则返回 false
func sleepOrDone(ctx context.Context, d time.Duration) bool {
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return false
	case <-timer.C:
		return true
	}
}

func (s *TaskQueueSubscriber) taskHandle(task *Task, handler TaskQueueHandler) {
	if task == nil {
		return
	}

	traceId := goo_utils.UUID()

	log := func() *goo_log.Entry {
		return s.log().WithTag("taskHandle").WithField("trace-id", traceId).WithField("task_id", task.Id)
	}

	percent, err := goo_utils.MemoryUsedPercent()
	if err != nil {
		log().WithField("task", task).ErrorF("获取任务，获取内存失败: %s", err.Error())
	} else {
		log().WithField("task", task).InfoF("获取任务，内存占用=%0.2f%%", percent)
	}

	// 任务执行不要挂在 Root() 上：进程退出信号只应停止拉新任务，在途任务需跑完/自行 Timeout
	ctx := goocontext.WithTraceId(context.Background(), traceId)
	ctx, cancel := goocontext.WithTimeout(ctx, time.Duration(task.Timeout)*time.Millisecond)
	defer cancel()

	// 续租跟随任务 ctx：Timeout 到期或 handler 返回后停止，避免卡死任务被无限续租
	renewCtx, renewCancel := context.WithCancel(ctx)
	defer renewCancel()
	goo_utils.AsyncFunc(func() {
		s.leaseRenewLoop(renewCtx, task, log)
	})

	startTime := time.Now()

	err = handler(ctx, task)
	renewCancel()

	if err == nil {
		log().WithField("执行时长_ms", time.Since(startTime).Milliseconds()).
			WithField("generation", task.Generation).
			InfoF("执行任务成功(%d/%d)", task.RetryTimes, task.MaxRetry)
		if err := s.TaskQueueTasks.taskDel(task); err != nil {
			if errors.Is(err, errStaleGeneration) {
				log().Warn("成功收尾跳过：generation 已过期")
			} else {
				log().WithField("task", task).ErrorF("成功收尾失败: %v", err)
			}
		}
		return
	}

	elapsedMs := time.Since(startTime).Milliseconds()
	reason := "执行任务失败"
	if errors.Is(err, context.DeadlineExceeded) {
		reason = "执行任务超时"
	}

	s.onFailure(task, err, log, elapsedMs, reason)
}

func (s *TaskQueueSubscriber) onFailure(task *Task, err error, log func() *goo_log.Entry, elapsedMs int64, reason string) {
	nextRunAtMs := time.Now().UnixMilli()
	if d, ok := retryAfterDuration(err); ok && d > 0 {
		nextRunAtMs = time.Now().Add(d).UnixMilli()
	}

	if task.MaxRetry != 0 && task.RetryTimes >= task.MaxRetry {
		log().WithField("执行时长_ms", elapsedMs).WithField("generation", task.Generation).
			WarnF("%s，达到最大重试次数(%d/%d)", reason, task.RetryTimes, task.MaxRetry)
		if failErr := s.TaskQueueTasks.taskFail(task); failErr != nil {
			if errors.Is(failErr, errStaleGeneration) {
				log().Warn("失败收尾跳过：generation 已过期")
			} else {
				log().WithField("task", task).ErrorF("失败收尾失败: %v", failErr)
			}
		}
		return
	}

	log().WithField("执行时长_ms", elapsedMs).WithField("next_run_at_ms", nextRunAtMs).
		WithField("generation", task.Generation).
		WarnF("%s，重试(%d/%d)", reason, task.RetryTimes, task.MaxRetry)
	if requeueErr := s.TaskQueueTasks.requeue(task, nextRunAtMs); requeueErr != nil {
		if errors.Is(requeueErr, errStaleGeneration) {
			log().Warn("重试收尾跳过：generation 已过期")
			return
		}
		log().WithField("task", task).ErrorF("重试收尾失败: %v", requeueErr)
		return
	}
	time.Sleep(time.Duration(rand.Intn(600)+200) * time.Millisecond)
}

// leaseRenewLoop 按 Timeout/3 刷新 processing 租约；间隔必须小于 Timeout，避免短任务被误回收
func (s *TaskQueueSubscriber) leaseRenewLoop(ctx context.Context, task *Task, log func() *goo_log.Entry) {
	timeout := time.Duration(task.Timeout) * time.Millisecond
	if timeout <= 0 {
		return
	}

	interval := timeout / 3
	if interval < 200*time.Millisecond {
		interval = 200 * time.Millisecond
	}
	if interval > 5*time.Minute {
		interval = 5 * time.Minute
	}
	// 续租周期必须短于租约窗口，否则 Leader 会在首次续租前误回收
	if interval >= timeout {
		interval = timeout / 2
		if interval < time.Millisecond {
			interval = time.Millisecond
		}
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := s.TaskQueueTasks.renewLease(task); err != nil {
				if errors.Is(err, errStaleGeneration) {
					log().Warn("续租停止：租约已失效（generation 过期或不在 processing）")
					return
				}
				log().WarnF("续租失败: %v", err)
			}
		}
	}
}

func (s *TaskQueueSubscriber) workId() string {
	localIp, _ := goo.LocalIP()
	return fmt.Sprintf("%s:%d", localIp, os.Getpid())
}

func (s *TaskQueueSubscriber) heartBeat(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			s.r.HSet(s.TaskWorkersKey, s.workId(), time.Now().Format("2006-01-02 15:04:05"))
			s.r.Expire(s.TaskWorkersKey, time.Second*10)

			if !sleepOrDone(ctx, time.Duration(rand.Intn(600)+200)*time.Millisecond) {
				return
			}
		}
	}
}

func (s *TaskQueueSubscriber) log() *goo_log.Entry {
	return goo_log.WithTag("goo-task-queue-subscribe", s.pid)
}
