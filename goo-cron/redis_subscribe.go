package goo_cron

import (
	"sync"

	"github.com/go-redis/redis"
	goo_context "github.com/liqiongtao/googo.io/goo-context"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"github.com/robfig/cron/v3"
)

var (
	cronTaskCode2EntryId sync.Map
)

func RedisSubscribe(r *redis.Client, c *cron.Cron, key string, tasks map[string]CronTaskFunc) {
	sub := r.Subscribe(key)
	defer func() { _ = sub.Close() }()

	for {
		select {
		case <-goo_context.WithCancel().Done():
			return

		case msg := <-sub.Channel():
			if msg.Channel != key {
				continue
			}

			task, err := ConvertCronTaskData(msg.Payload)
			if err != nil {
				goo_log.WithField("payload", msg.Payload).ErrorF("convert cron task data err: %v", err)
				continue
			}
			if err = task.Valid(); err != nil {
				goo_log.WithField("task", task).ErrorF("validate cron task data err: %v", err)
				continue
			}

			taskFunc, ok := tasks[task.Code]
			if !ok {
				goo_log.WithField("task", task).ErrorF("invalid cron task code: %s", task.Code)
				continue
			}

			switch task.Status {
			case CronTaskStatusDelete:
				if entryId, ok := cronTaskCode2EntryId.Load(task.Code); ok {
					cronTaskCode2EntryId.Delete(task.Code)
					c.Remove(entryId.(cron.EntryID))
				}

			case CronTaskStatusCreate:
				entryId, err := c.AddFunc(task.Spec, taskFunc(task))
				if err != nil {
					goo_log.WithField("task", task).ErrorF("add cron task err: %v", err)
					continue
				}
				
				cronTaskCode2EntryId.Store(task.Code, entryId)

			case CronTaskStatusUpdate:
				if entryId, ok := cronTaskCode2EntryId.Load(task.Code); ok {
					c.Remove(entryId.(cron.EntryID))
				}

				entryId, err := c.AddFunc(task.Spec, taskFunc(task))
				if err != nil {
					goo_log.WithField("task", task).ErrorF("add cron task err: %v", err)
					continue
				}

				cronTaskCode2EntryId.Store(task.Code, entryId)
			}
		}
	}
}
