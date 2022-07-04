package goo

import (
	goo_redis "github.com/liqiongtao/googo.io/goo-redis"
)

func Redis(names ...string) *goo_redis.Client {
	return goo_redis.GetClient(names...)
}
