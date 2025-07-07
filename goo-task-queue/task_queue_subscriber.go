package goo_task_queue

import (
	"context"
	"errors"
	"github.com/go-redis/redis"
	goo_context "github.com/liqiongtao/googo.io/goo-context"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
	"runtime"
	"time"
)

type TaskQueueHandler func(ctx context.Context, value any) error

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

	// 并发数控制
	if limit <= 0 {
		limit = runtime.NumCPU() * 2
	}

	s.log().InfoF("任务监听成功，并发数: %d", limit)

	var (
		limitCH = make(chan any, limit)
		taskCH  = make(chan *Task, limit)

		done      = make(chan any)
		cancelCtx = goo_context.WithCancel()
	)

	// 执行任务
	goo_utils.AsyncFunc(func() {
		defer func() { done <- struct{}{} }()

		for {
			select {
			case <-cancelCtx.Done():
				s.log().Info("任务执行协程退出")
				return

			case task, ok := <-taskCH:
				if !ok {
					s.log().Warn("未获取到任务")
					return
				}

				// 执行任务
				s.taskHandle(task, handler)

				<-limitCH
			}
		}
	})

	// 获取任务
	goo_utils.AsyncFunc(func() {
		defer func() { done <- struct{}{} }()

		for {
			select {
			case <-cancelCtx.Done():
				s.log().Info("获取任务协程退出")
				return

			case limitCH <- struct{}{}:
				goo_utils.AsyncFunc(func() {
					task, err := s.getOneTask()
					if err != nil || task == nil {
						time.Sleep(time.Second * 3)
						<-limitCH
						return
					}
					taskCH <- task
				})
			}
		}
	})

	for i := 0; i < 2; i++ {
		<-done
	}

	close(limitCH)
	close(taskCH)
	close(done)
}

func (s *TaskQueueSubscriber) taskHandle(task *Task, handler TaskQueueHandler) {
	// 追踪ID
	traceId := goo_utils.UUID()

	// 任务日志
	log := func() *goo_log.Entry {
		return s.log().WithTag("taskHandle").WithField("trace-id", traceId)
	}
	log().WithField("task", task).Info("获取任务")

	ctx := context.Background()
	ctx = context.WithValue(ctx, "trace-id", traceId)

	// 执行成功
	if err := handler(ctx, task); err == nil {
		log().Info("任务执行成功")
		s.TaskQueueTasks.taskDel(task.Id)
		return
	}

	// 达到最大执行次数
	if task.MaxRetry != 0 && task.RetryTimes >= task.MaxRetry {
		log().Warn("任务执行失败")
		s.taskFail(task)
		return
	}

	log().Warn("任务执行失败，重试")

	// 增加重试次数
	s.retry(task)
}

// 重试
func (s *TaskQueueSubscriber) retry(tasks ...*Task) error {
	pi := s.r.TxPipeline()

	for _, task := range tasks {
		task.RetryTimes++
		// 重试次数+1
		pi.HIncrBy(s.taskInfoKey(task.Id), "retry_times", 1)
		// 删除执行队列
		pi.ZRem(s.TaskProcessingKey, task.Id)
		// 添加待执行队列
		pi.ZAdd(s.TaskPendingKey, redis.Z{Member: task.Id, Score: float64(time.Now().Unix() + 60)})
	}

	if _, err := pi.Exec(); err != nil {
		s.log().WithTag("retry").Error(err)
		return err
	}

	return nil
}

// 任务失败
func (s *TaskQueueSubscriber) taskFail(tasks ...*Task) error {
	pi := s.r.TxPipeline()

	for _, task := range tasks {
		// 删除执行队列
		pi.ZRem(s.TaskProcessingKey, task.Id)
		// 添加失败队列
		pi.ZAdd(s.TaskFailKey, redis.Z{Member: task.Id, Score: float64(time.Now().Unix())})
	}

	if _, err := pi.Exec(); err != nil {
		s.log().WithTag("taskFail").Error(err)
		return err
	}

	return nil
}

func (s *TaskQueueSubscriber) log() *goo_log.Entry {
	return s.TaskQueue.log().WithTag("goo-task-queue-subscribe")
}
