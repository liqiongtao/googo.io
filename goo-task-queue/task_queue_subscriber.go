package goo_task_queue

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/liqiongtao/googo.io/goo"
	goo_context "github.com/liqiongtao/googo.io/goo-context"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
	"os"
	"runtime"
	"time"
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
		limitCH = make(chan any, limit)
		taskCH  = make(chan *Task, limit)

		done = make(chan any)
	)

	// 执行任务
	goo_utils.AsyncFunc(func() {
		defer func() { done <- struct{}{} }()

		var (
			canExit bool
		)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		for {
			select {
			case <-ctx.Done():
				return

			case <-goo_context.WithCancel().Done():
				canExit = true
				if len(taskCH) == 0 {
					s.log().Info("任务执行为空，任务执行退出")
					return
				}

			case task, ok := <-taskCH:
				if !ok {
					cancel()
					s.log().Warn("未获取到任务")
					return
				}

				goo_utils.AsyncFunc(func() {
					defer func() {
						<-limitCH

						if canExit && len(limitCH) == 0 {
							cancel()
							s.log().Info("任务执行完毕，任务执行退出")
							return
						}
					}()

					// 执行任务
					s.taskHandle(task, handler)
				})
			}
		}
	})

	// 获取任务
	goo_utils.AsyncFunc(func() {
		defer func() { done <- struct{}{} }()

		for {
			select {
			case <-goo_context.WithCancel().Done():
				s.log().Info("获取任务协程退出")
				return

			case limitCH <- struct{}{}:
				goo_utils.AsyncFunc(func() {
					task, err := s.getOneTask()
					if err != nil || task == nil {
						time.Sleep(time.Second)
						<-limitCH
						return
					}
					taskCH <- task
				})
			}
		}
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
		close(taskCH)
		close(limitCH)
		close(done)
	})
}

func (s *TaskQueueSubscriber) taskHandle(task *Task, handler TaskQueueHandler) {
	// 追踪ID
	traceId := goo_utils.UUID()

	// 任务日志
	log := func() *goo_log.Entry {
		return s.log().WithTag("taskHandle").WithField("trace-id", traceId).WithField("task_id", task.Id)
	}
	log().WithField("task", task).Info("获取任务")

	// 上下文
	ctx := context.Background()
	ctx = context.WithValue(ctx, "trace-id", traceId)

	// 超时控制
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(task.Timeout))
	defer cancel()

	// 任务执行
	err := handler(ctx, task)

	// 执行成功
	if err == nil {
		log().Info("执行任务成功")
		s.TaskQueueTasks.taskDel(task.Id)
		return
	}

	// 任务执行超时
	if errors.Is(err, context.DeadlineExceeded) {
		log().Warn("执行任务超时, 准备重试")
		s.retry(task)
		return
	}

	// 达到最大执行次数
	if task.MaxRetry != 0 && task.RetryTimes >= task.MaxRetry {
		log().Warn("执行任务失败，达到最大重试次数")
		s.taskFail(task)
		return
	}

	log().Warn("执行任务失败，重试")

	// 增加重试次数
	s.retry(task)

	time.Sleep(time.Second)
}

// 重试
func (s *TaskQueueSubscriber) retry(tasks ...*Task) error {
	pi := s.r.TxPipeline()

	for _, task := range tasks {
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

			time.Sleep(time.Second)
		}
	}
}

func (s *TaskQueueSubscriber) log() *goo_log.Entry {
	return goo_log.WithTag("goo-task-queue-subscribe")
}
