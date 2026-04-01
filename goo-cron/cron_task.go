package goo_cron

import (
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
	code2Func    map[string]TaskFunc
}

func New(key string, code2Func map[string]TaskFunc, opts ...Option) *CronTask {
	c := &CronTask{
		c:         cron.New(cron.WithSeconds()),
		key:       key,
		code2Func: code2Func,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *CronTask) Cron() *cron.Cron {
	return c.c
}

func (c *CronTask) Run() {
	ctx := goo_context.WithCancel()

	goo_utils.AsyncFunc(func() {
		if c.r == nil {
			return
		}
		c.Subscribe(ctx)
	})

	c.c.Start()

	<-ctx.Done()
	goo_log.WithTag("goo-cron").Debug("系统退出，等待全部任务执行结束...")

	for _, entry := range c.c.Entries() {
		c.c.Remove(entry.ID)
	}

	time.Sleep(time.Second)

	<-c.c.Stop().Done()
	goo_log.WithTag("goo-cron").Debug("系统退出成功，全部任务执行结束")
}

func (c *CronTask) Add(task *TaskData, hooks ...TaskFunc) error {
	entryId, err := c.c.AddFunc(task.Spec, func() {
		handler, ok := c.code2Func[task.Code]
		if !ok {
			goo_log.WithTag("goo-cron").WithField("task", task).Warn("no task handler")
			return
		}

		defer func() {
			for _, hook := range hooks {
				hook(task)
			}
		}()

		handler(task)
	})

	if err == nil {
		c.code2EntryId.Store(task.Code, entryId)
	}

	return err
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

func (c *CronTask) Subscribe(ctx *goo_context.Context) {
	sub := c.r.Subscribe(c.key)
	defer func() { _ = sub.Close() }()

	for {
		select {
		case <-ctx.Done():
			goo_log.WithTag("goo-cron").Info("定时任务订阅服务退出")
			return

		case msg := <-sub.Channel():
			if msg.Channel != c.key {
				continue
			}
			if msg.Payload == "" {
				goo_log.WithTag("goo-cron").WithField("msg", msg).Warn("payload is empty")
				continue
			}

			task, err := ConvertTaskData(msg.Payload)
			if err != nil {
				goo_log.WithTag("goo-cron").WithField("payload", msg.Payload).ErrorF("convert cron task data err: %v", err)
				continue
			}
			if err = task.Valid(); err != nil {
				goo_log.WithTag("goo-cron").WithField("task", task).ErrorF("validate cron task data err: %v", err)
				continue
			}

			goo_log.WithTag("goo-cron").WithField("task", task).Info("receive message")

			switch task.Status {
			case TaskStatusDelete: // 删除任务
				c.Remove(task.Code)

			case TaskStatusCreate, TaskStatusUpdate: // 添加、更新任务
				c.Remove(task.Code)
				if err := c.Add(task); err != nil {
					goo_log.WithTag("goo-cron").WithField("task", task).ErrorF("add cron task err: %v", err)
					continue
				}

			case TaskStatusExecute: // 立即执行
				handler, ok := c.code2Func[task.Code]
				if !ok {
					goo_log.WithTag("goo-cron").WithField("task", task).Warn("no task handler")
					continue
				}
				goo_utils.AsyncFunc(func() {
					handler(task)
				})
			}
		}
	}
}
