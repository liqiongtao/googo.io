package goo_cron

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis"
	goo_context "github.com/liqiongtao/googo.io/goo-context"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
	"github.com/robfig/cron/v3"
)

type CronTask struct {
	r            *redis.Client
	c            *cron.Cron
	key          string
	code2EntryId sync.Map
	jobs         map[string]TaskFunc
}

func New(key string, r *redis.Client, jobs map[string]TaskFunc, opts ...cron.Option) *CronTask {
	opts = append(opts, cron.WithSeconds())
	return &CronTask{
		r:    r,
		c:    cron.New(opts...),
		key:  key,
		jobs: jobs,
	}
}

func (c *CronTask) Cron() *cron.Cron {
	return c.c
}

func (c *CronTask) Run() {
	goo_utils.AsyncFunc(c.Subscribe)

	c.c.Start()

	<-goo_context.WithCancel().Done()
	goo_log.WithTag("goo-cron").Debug("系统退出，等待全部任务执行结束...")

	for _, entry := range c.c.Entries() {
		c.c.Remove(entry.ID)
	}

	time.Sleep(time.Second)

	<-c.c.Stop().Done()
	goo_log.WithTag("goo-cron").Debug("系统退出成功，全部任务执行结束")
}

func (c *CronTask) Add(task *TaskData) error {
	taskFunc, ok := c.jobs[task.Code]
	if !ok {
		return fmt.Errorf("task code %s not exists", task.Code)
	}

	entryId, err := c.c.AddFunc(task.Spec, taskFunc(task))
	if err != nil {
		return fmt.Errorf("add task %s err: %v", task.Code, err)
	}

	c.code2EntryId.Store(task.Code, entryId)

	return nil
}

func (c *CronTask) Remove(taskCode string) {
	v, ok := c.code2EntryId.Load(taskCode)
	if ok {
		c.code2EntryId.Delete(taskCode)
	}
	if entryId, ok := v.(cron.EntryID); ok {
		c.c.Remove(entryId)
	}
}

func (c *CronTask) Subscribe() {
	sub := c.r.Subscribe(c.key)
	defer func() { _ = sub.Close() }()

	for {
		select {
		case <-goo_context.WithCancel().Done():
			goo_log.Info("定时任务订阅服务退出")
			return

		case msg := <-sub.Channel():
			if msg.Channel != c.key {
				continue
			}

			task, err := ConvertTaskData(msg.Payload)
			if err != nil {
				goo_log.WithField("payload", msg.Payload).ErrorF("convert cron task data err: %v", err)
				continue
			}
			if err = task.Valid(); err != nil {
				goo_log.WithField("task", task).ErrorF("validate cron task data err: %v", err)
				continue
			}

			goo_log.WithField("task", task).Info("receive message")

			switch task.Status {
			case TaskStatusDelete:
				c.Remove(task.Code)

			case TaskStatusCreate:
				if err := c.Add(task); err != nil {
					goo_log.WithField("task", task).ErrorF("add cron task err: %v", err)
					continue
				}

			case TaskStatusUpdate:
				c.Remove(task.Code)
				if err := c.Add(task); err != nil {
					goo_log.WithField("task", task).ErrorF("add cron task err: %v", err)
					continue
				}
			}
		}
	}
}
