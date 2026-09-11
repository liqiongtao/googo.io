package goo_kafka

import goo_redis "github.com/liqiongtao/googo.io/goo-redis"

const (
	FocusName = "focus"
	RedisName = "redis"
)

type Option struct {
	Name  string
	Value any
}

// 是否强制
func FocusOption() Option {
	return Option{Name: FocusName, Value: true}
}

// redis
func RedisOption(cli *goo_redis.Client) Option {
	return Option{Name: RedisName, Value: cli}
}
