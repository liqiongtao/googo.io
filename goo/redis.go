package goo

import (
	gooredis "github.com/liqiongtao/googo.io/goo-redis"
)

func Redis(names ...string) *gooredis.Client {
	return gooredis.GetClient(names...)
}
