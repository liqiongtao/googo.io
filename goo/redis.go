package goo

import (
	goo_redis "github.com/liqiongtao/googo.io/goo-redis"
)

func Redis(names ...string) *goo_redis.Redis {
	return goo_redis.Client(names...)
}
