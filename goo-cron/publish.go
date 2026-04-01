package goo_cron

import (
	"github.com/go-redis/redis"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
)

func Publish(r *redis.Client, key string, task *TaskData) error {
	payload := task.String()
	goo_log.WithTag("goo-cron").WithField("payload", payload).Info("发布任务")
	return r.Publish(key, payload).Err()
}
