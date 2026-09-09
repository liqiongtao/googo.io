package goo_cron

import (
	goo_redis "github.com/liqiongtao/googo.io/goo-redis"
)

type Option func(*CronTask)

func WithRedis(r *goo_redis.Client) Option {
	return func(c *CronTask) {
		c.r = r
	}
}

func WithHooks(hooks ...TaskFunc) Option {
	return func(c *CronTask) {
		c.hooks = append(c.hooks, hooks...)
	}
}
