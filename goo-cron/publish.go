package goo_cron

import (
	"errors"

	goolog "github.com/liqiongtao/googo.io/goo-log"
	gooredis "github.com/liqiongtao/googo.io/goo-redis"
)

func Publish(r *gooredis.Client, key string, task *TaskData) error {
	if r == nil {
		goolog.WithTag("goo-cron").Error("publish cron redis client is nil")
		return errors.New("cron redis client is nil")
	}
	if task == nil {
		goolog.WithTag("goo-cron").Error("publish cron task is nil")
		return errors.New("cron task is nil")
	}

	if err := task.Valid(); err != nil {
		goolog.WithTag("goo-cron").WithField("task", task).ErrorF("validate cron task data err: %v", err)
		return err
	}

	err := r.Publish(key, task.String()).Err()
	if err != nil {
		goolog.WithTag("goo-cron").WithField("task", task).ErrorF("publish cron task err: %v", err)
	} else {
		goolog.WithTag("goo-cron").WithField("task", task).Info("publish cron task success")
	}
	return err
}
