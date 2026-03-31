package goo_cron

import (
	"github.com/go-redis/redis"
)

func Publish(r *redis.Client, key string, task TaskData) error {
	return r.Publish(key, task.String()).Err()
}
