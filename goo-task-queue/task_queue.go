package goo_task_queue

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	goo_context "github.com/liqiongtao/googo.io/goo-context"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	goo_redis "github.com/liqiongtao/googo.io/goo-redis"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
	"runtime"
	"time"
)

type TaskQueue struct {
	r *goo_redis.Client

	taskInfoKey       string // 任务信息
	taskPendingKey    string // 待处理任务
	taskProcessingKey string // 正在处理任务
	taskFailKey       string // 失败任务
}

func New(r *goo_redis.Client) *TaskQueue {
	t := &TaskQueue{
		r:                 r,
		taskInfoKey:       TaskInfoKey,
		taskPendingKey:    TaskPendingKey,
		taskProcessingKey: TaskProcessingKey,
		taskFailKey:       TaskFailKey,
	}
	if r.Config.Prefix != "" {
		t.WithPrefix(r.Config.Prefix)
	}
	return t
}

func (t *TaskQueue) WithPrefix(prefix string) *TaskQueue {
	t.taskInfoKey = fmt.Sprintf("%s:%s", prefix, TaskInfoKey)
	t.taskPendingKey = fmt.Sprintf("%s:%s", prefix, TaskPendingKey)
	t.taskProcessingKey = fmt.Sprintf("%s:%s", prefix, TaskProcessingKey)
	t.taskFailKey = fmt.Sprintf("%s:%s", prefix, TaskFailKey)
	return t
}

func (t *TaskQueue) Publish(tasks ...*Task) error {
	log := func() *goo_log.Entry {
		return t.log().WithTag("Publish")
	}

	if t.r == nil {
		err := errors.New("redis未初始化")
		log().Error(err)
		return err
	}

	pi := t.r.TxPipeline()

	for _, task := range tasks {
		if task.MaxRetry == 0 {
			task.MaxRetry = 99 // 默认重试次数 99次
		}
		if task.RetryIntervalTime == 0 {
			task.RetryIntervalTime = 1 // 默认重试间隔时间 1秒
		}
		if task.Timeout == 0 {
			task.Timeout = 1800 // 默认超时时间 30分钟
		}

		pi.HMSet(fmt.Sprintf("%s:%s", t.taskInfoKey, task.Id), task.MapData())
		pi.ZAdd(t.taskPendingKey, redis.Z{Score: task.Score(), Member: task.Id})
	}

	if _, err := pi.Exec(); err != nil {
		log().WithField("tasks", tasks).ErrorF("发布任务失败: %s", err.Error())
		return err
	}

	log().WithField("tasks", tasks).Info("发布任务成功")

	return nil
}

func (t *TaskQueue) Subscribe(limit int, handler TaskQueueHandler) {
	log := func() *goo_log.Entry {
		return t.log().WithTag("Subscribe")
	}

	if t.r == nil {
		err := errors.New("redis未初始化")
		log().Error(err)
		return
	}

	// 并发数控制
	if limit <= 0 {
		limit = runtime.NumCPU() * 2
	}

	log().InfoF("任务监听成功，并发数: %d", limit)

	var (
		limitCH = make(chan any, limit)
		taskCH  = make(chan *Task, limit)

		done   = make(chan any)
		cancel = goo_context.WithCancel()
	)

	// 执行任务
	goo_utils.AsyncFunc(func() {
		defer func() { done <- struct{}{} }()

		for {
			select {
			case <-cancel.Done():
				log().Info("任务执行协程退出")
				return

			case task, ok := <-taskCH:
				if !ok {
					log().Warn("未获取到任务")
					return
				}

				// 追踪ID
				traceId := goo_utils.UUID()

				// 任务日志
				taskLog := func() *goo_log.Entry {
					return log().WithField("trace-id", traceId)
				}
				taskLog().WithField("task", task).Info("获取任务")

				ctx := context.Background()
				ctx = context.WithValue(ctx, "trace-id", traceId)

				// 执行成功
				if err := handler(ctx, task); err == nil {
					taskLog().Info("任务执行成功")
					t.delTask(task)
					<-limitCH
					continue
				}

				// 达到最大执行次数
				if task.MaxRetry != 0 && task.RetryTimes >= task.MaxRetry {
					taskLog().Warn("任务执行失败")
					t.taskFail(task)
					<-limitCH
					continue
				}

				taskLog().Warn("任务执行失败，重试")

				// 增加重试次数
				t.retry(task)

				// 暂停
				time.Sleep(time.Second * time.Duration(task.RetryIntervalTime))

				<-limitCH
			}
		}
	})

	// 获取任务
	goo_utils.AsyncFunc(func() {
		defer func() { done <- struct{}{} }()

		for {
			select {
			case <-cancel.Done():
				log().Info("获取任务协程退出")
				return

			case limitCH <- struct{}{}:
				task, err := t.getOneTask()
				if err != nil {
					time.Sleep(time.Second * 3)
					<-limitCH
					continue
				}
				if task == nil {
					time.Sleep(time.Second)
					<-limitCH
					continue
				}
				taskCH <- task
			}
		}
	})

	for i := 0; i < 2; i++ {
		<-done
	}

	close(limitCH)
	close(taskCH)
}

func (t *TaskQueue) getOneTask() (*Task, error) {
	// 获取任务
	z := t.r.BZPopMin(time.Second*10, t.taskPendingKey).Val()
	if z.Member == nil {
		t.log().WithTag("getOneTask").Warn("member is nil")
		return nil, nil
	}

	// 添加执行队列
	if err := t.r.ZAdd(t.taskProcessingKey, redis.Z{Member: z.Member, Score: z.Score}).Err(); err != nil {
		t.log().WithTag("getOneTask").Error("add taskProcessingKey fail", err)
		t.r.ZAdd(t.taskPendingKey, redis.Z{Member: z.Member, Score: z.Score})
		return nil, err
	}

	// 任务ID
	taskId := z.Member.(string)

	// 查询任务
	task, err := t.getTaskById(taskId)
	if err != nil {
		t.log().WithTag("getOneTask").Error("get task fail", err)
		t.r.ZAdd(t.taskPendingKey, redis.Z{Member: z.Member, Score: z.Score})
		return nil, err
	}

	return task, nil
}

// 根据ID获取任务
func (t *TaskQueue) getTaskById(taskId string) (*Task, error) {
	// 任务KEY
	taskKey := fmt.Sprintf("%s:%s", t.taskInfoKey, taskId)

	// 如果key不存在
	if t.r.Exists(taskKey).Val() == 0 {
		// 删除执行队列
		t.r.ZRem(t.taskProcessingKey, taskId)
		// 删除待执行队列
		t.r.ZRem(t.taskPendingKey, taskId)

		err := fmt.Errorf("task %s not found", taskId)
		t.log().WithTag("getTaskById").Error(err)
		return nil, err
	}

	task := &Task{
		Id:      t.r.HGet(taskKey, "id").Val(),
		Type:    t.r.HGet(taskKey, "type").Val(),
		Payload: t.r.HGet(taskKey, "payload").Val(),
	}

	task.HighPriority, _ = t.r.HGet(taskKey, "high_priority").Int()
	task.MaxRetry, _ = t.r.HGet(taskKey, "max_retry").Int()
	task.RetryTimes, _ = t.r.HGet(taskKey, "retry_times").Int()
	task.RetryIntervalTime, _ = t.r.HGet(taskKey, "retry_interval_time").Int()
	task.Timeout, _ = t.r.HGet(taskKey, "timeout").Int64()

	return task, nil
}

// 重试
func (t *TaskQueue) retry(tasks ...*Task) error {
	pi := t.r.TxPipeline()

	for _, task := range tasks {
		task.RetryTimes++
		// 重试次数+1
		pi.HIncrBy(fmt.Sprintf("%s:%s", t.taskInfoKey, task.Id), "retry_times", 1)
		// 删除执行队列
		pi.ZRem(t.taskProcessingKey, task.Id)
		// 添加待执行队列
		pi.ZAdd(t.taskPendingKey, redis.Z{Member: task.Id, Score: task.Score()})
	}

	if _, err := pi.Exec(); err != nil {
		t.log().WithTag("retry").Error(err)
		return err
	}

	return nil
}

// 任务失败
func (t *TaskQueue) taskFail(tasks ...*Task) error {
	pi := t.r.TxPipeline()

	for _, task := range tasks {
		// 删除执行队列
		pi.ZRem(t.taskProcessingKey, task.Id)
		// 添加失败队列
		pi.ZAdd(t.taskFailKey, redis.Z{Member: task.Id, Score: task.Score()})
	}

	if _, err := pi.Exec(); err != nil {
		t.log().WithTag("taskFail").Error(err)
		return err
	}

	return nil
}

// 删除任务
func (t *TaskQueue) delTask(tasks ...*Task) error {
	pi := t.r.TxPipeline()

	for _, task := range tasks {
		// 删除任务
		pi.Del(fmt.Sprintf("%s:%s", t.taskInfoKey, task.Id))
		// 删除执行队列
		pi.ZRem(t.taskProcessingKey, task.Id)
	}

	if _, err := pi.Exec(); err != nil {
		t.log().WithTag("delTask").Error(err)
		return err
	}

	return nil
}

func (t *TaskQueue) log() *goo_log.Entry {
	return goo_log.WithTag("goo-task-queue")
}
