package goo_cron

import (
	gooredis "github.com/liqiongtao/googo.io/goo-redis"
)

type Option func(*CronTask)

func WithRedis(r *gooredis.Client) Option {
	return func(c *CronTask) {
		c.r = r
	}
}

func WithHooks(hooks ...TaskFunc) Option {
	return func(c *CronTask) {
		c.hooks = append(c.hooks, hooks...)
	}
}
