package goo_cron

import (
	"github.com/go-redis/redis"
)

func RedisPublish(key string, r *redis.Client, d CronTaskData) error {
	return r.Publish(key, d.String()).Err()
}
