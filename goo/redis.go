package goo

import (
	goo_redis "googo.io/goo-redis"
)

func Redis() *goo_redis.Redis {
	return goo_redis.Client()
}
