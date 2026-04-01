package goo_cron

import (
	"github.com/go-redis/redis"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
)

func Publish(r *redis.Client, key string, task *TaskData) error {
	if err := task.Valid(); err != nil {
		goo_log.WithTag("goo-cron").WithField("task", task).ErrorF("validate cron task data err: %v", err)
		return err
	}

	err := r.Publish(key, task.String()).Err()
	if err != nil {
		goo_log.WithTag("goo-cron").WithField("task", task).ErrorF("publish cron task err: %v", err)
	} else {
		goo_log.WithTag("goo-cron").WithField("task", task).Info("publish cron task success")
	}
	return err
}
