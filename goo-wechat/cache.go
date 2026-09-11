package goowechat

import (
	gooredis "github.com/liqiongtao/googo.io/goo-redis"
)

var (
	__cache *gooredis.Client
)

func InitCache(redisClient *gooredis.Client) {
	__cache = redisClient
}
