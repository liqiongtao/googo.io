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
	"math/rand"
	"os"
	"runtime"
	"sync"
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

		var wg sync.WaitGroup

		for {
			task, ok := <-taskCH
			if !ok {
				break
			}

			wg.Add(1)

			goo_utils.AsyncFunc(func() {
				defer func() {
					<-limitCH
					wg.Done()
				}()

				// 执行任务
				s.taskHandle(task, handler)
			})
		}

		wg.Wait()

		s.log().Info("任务执行完毕，任务执行退出")
	})

	// 获取任务
	goo_utils.AsyncFunc(func() {
		defer func() { done <- struct{}{} }()

		var wg sync.WaitGroup

		for {
			// 1. 检查是否退出
			if willExit {
				close(limitCH)
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

			// 4. 执行任务
			goo_utils.AsyncFunc(func() {
				defer wg.Done()

				// 获取任务
				task, err := s.getOneTask()
				if err != nil || task == nil {
					time.Sleep(time.Duration(rand.Intn(600)+200) * time.Millisecond)
					<-limitCH
					return
				}

				// 发布任务
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

	// 超时控制
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(task.Timeout))
	defer cancel()

	// 任务执行
	err = handler(ctx, task)

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

	time.Sleep(time.Duration(rand.Intn(600)+200) * time.Millisecond)
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
		pi.ZAdd(s.TaskPendingKey, redis.Z{Member: task.Id, Score: float64(time.Now().Unix()) + float64(rand.Intn(60)+60)})
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

			time.Sleep(time.Duration(rand.Intn(600)+200) * time.Millisecond)
		}
	}
}

func (s *TaskQueueSubscriber) log() *goo_log.Entry {
	return goo_log.WithTag("goo-task-queue-subscribe")
}
