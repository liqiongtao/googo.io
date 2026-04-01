package goo_cron

import (
	"github.com/go-redis/redis"
)

type Option func(*CronTask)

func WithRedis(r *redis.Client) Option {
	return func(c *CronTask) {
		c.r = r
	}
}

func WithHooks(hooks ...TaskFunc) Option {
	return func(c *CronTask) {
		c.hooks = append(c.hooks, hooks...)
	}
}
