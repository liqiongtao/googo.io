package goo_task_queue

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/liqiongtao/googo.io/goo"
	goo_context "github.com/liqiongtao/googo.io/goo-context"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
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

	// 选举leader
	goo_utils.AsyncFunc(func() {
		s.TaskQueueLeader.Generate()
	})

	// 心跳
	goo_utils.AsyncFunc(func() {
		s.heartBeat()
	})

	// 写入pid文件
	goo_utils.AsyncFunc(func() {
		goo_utils.WriteToFile(".pid", []byte(fmt.Sprintf("%d", os.Getpid())))
	})

	// 并发数控制
	if limit <= 0 {
		limit = runtime.NumCPU() * 2
	}

	s.log().InfoF("任务监听成功 并发数=%d workerId=%s", limit, s.workId())

	var (
		willExit bool

		limitCH = make(chan any, limit)
		taskCH  = make(chan *Task, limit)

		done = make(chan any)
	)

	// 监听退出
	goo_utils.AsyncFunc(func() {
		for {
			select {
			case <-goo_context.WithCancel().Done():
				willExit = true
				s.log().Info("准备退出")
				return
			}
		}
	})

	// 执行任务
	goo_utils.AsyncFunc(func() {
		defer func() { done <- struct{}{} }()

		var (
			wg  sync.WaitGroup
			cnt int
		)

		for {
			task, ok := <-taskCH
			if !ok {
				break
			}

			cnt++
			wg.Add(1)

			goo_utils.AsyncFunc(func() {
				defer func() {
					<-limitCH
					cnt--
					wg.Done()
					s.log().InfoF("剩余 %d 任务正在执行", cnt)
				}()

				// 执行任务
				s.taskHandle(task, handler)
			})
		}

		wg.Wait()

		s.log().Info("全部任务执行完毕，执行退出")
	})

	// 获取任务
	goo_utils.AsyncFunc(func() {
		defer func() { done <- struct{}{} }()

		var wg sync.WaitGroup

		for {
			// 1. 检查是否退出：不再向 limitCH 申请槽位，也不 close（避免 send on closed channel）
			if willExit {
				break
			}

			// 2. 检查内存
			percent, err := goo_utils.MemoryUsedPercent()
			if err != nil || percent >= s.MaxMemoryPercent {
				n := rand.Intn(600) + 3000
				s.log().WarnF("内存占用超过最大限制=%0.2f%%，等待%dms后重试", percent, n)
				time.Sleep(time.Duration(n) * time.Millisecond)
				continue
			}

			// 3. 并发控制
			wg.Add(1)
			limitCH <- struct{}{}

			// 4. 拉取任务
			goo_utils.AsyncFunc(func() {
				defer wg.Done()

				task, err := s.getOneTask()
				if err != nil || task == nil {
					time.Sleep(time.Duration(rand.Intn(600)+200) * time.Millisecond)
					<-limitCH
					return
				}

				taskCH <- task
			})
		}

		wg.Wait()
		close(taskCH)
	})

	// 监听退出信号
	for i := 0; i < 2; i++ {
		<-done
	}

	goo_utils.AsyncFuncGroup(func() {
		// 删除节点
		s.r.HDel(s.TaskWorkersKey, s.workId())
	}, func() {
		// 关闭管道
		close(done)
	})
}

func (s *TaskQueueSubscriber) taskHandle(task *Task, handler TaskQueueHandler) {
	if task == nil {
		return
	}

	// 追踪ID
	traceId := goo_utils.UUID()

	// 任务日志
	log := func() *goo_log.Entry {
		return s.log().WithTag("taskHandle").WithField("trace-id", traceId).WithField("task_id", task.Id)
	}

	percent, err := goo_utils.MemoryUsedPercent()
	if err != nil {
		log().WithField("task", task).ErrorF("获取任务，获取内存失败: %s", err.Error())
	} else {
		log().WithField("task", task).InfoF("获取任务，内存占用=%0.2f%%", percent)
	}

	// 上下文
	ctx := context.Background()
	ctx = context.WithValue(ctx, "trace-id", traceId)

	// 超时控制（Timeout 单位：毫秒）
	ctx, cancel := context.WithTimeout(ctx, time.Duration(task.Timeout)*time.Millisecond)
	defer cancel()

	startTime := time.Now()

	// 任务执行
	err = handler(ctx, task)

	// 执行成功
	if err == nil {
		log().WithField("执行时长_ms", time.Since(startTime).Milliseconds()).InfoF("执行任务成功(%d/%d)", task.RetryTimes, task.MaxRetry)
		s.TaskQueueTasks.taskDel(task.Id)
		return
	}

	elapsedMs := time.Since(startTime).Milliseconds()
	reason := "执行任务失败"
	if errors.Is(err, context.DeadlineExceeded) {
		reason = "执行任务超时"
	}

	// 超时与普通失败统一：达到 MaxRetry 进 fail，否则重试（业务可通过 RetryAfter 控制延迟）
	s.onFailure(task, err, log, elapsedMs, reason)
}

// onFailure 统一处理失败（含超时）
func (s *TaskQueueSubscriber) onFailure(task *Task, err error, log func() *goo_log.Entry, elapsedMs int64, reason string) {
	nextRunAtMs := time.Now().UnixMilli()
	if d, ok := retryAfterDuration(err); ok && d > 0 {
		nextRunAtMs = time.Now().Add(d).UnixMilli()
	}

	if task.MaxRetry != 0 && task.RetryTimes >= task.MaxRetry {
		log().WithField("执行时长_ms", elapsedMs).WarnF("%s，达到最大重试次数(%d/%d)", reason, task.RetryTimes, task.MaxRetry)
		_ = s.TaskQueueTasks.taskFail(task)
		return
	}

	log().WithField("执行时长_ms", elapsedMs).WithField("next_run_at_ms", nextRunAtMs).
		WarnF("%s，重试(%d/%d)", reason, task.RetryTimes, task.MaxRetry)
	_ = s.TaskQueueTasks.requeue(task, nextRunAtMs)
	time.Sleep(time.Duration(rand.Intn(600)+200) * time.Millisecond)
}

// 节点ID
func (s *TaskQueueSubscriber) workId() string {
	localIp, _ := goo.LocalIP()
	return fmt.Sprintf("%s:%d", localIp, os.Getpid())
}

// 节点心跳
func (s *TaskQueueSubscriber) heartBeat() {
	for {
		select {
		case <-goo_context.WithCancel().Done():
			return

		default:
			s.r.HSet(s.TaskWorkersKey, s.workId(), time.Now().Format("2006-01-02 15:04:05"))
			s.r.Expire(s.TaskWorkersKey, time.Second*10)

			time.Sleep(time.Duration(rand.Intn(600)+200) * time.Millisecond)
		}
	}
}

func (s *TaskQueueSubscriber) log() *goo_log.Entry {
	return goo_log.WithTag("goo-task-queue-subscribe", s.pid)
}
