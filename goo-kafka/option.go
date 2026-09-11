package gookafka

import gooredis "github.com/liqiongtao/googo.io/goo-redis"

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
func RedisOption(cli *gooredis.Client) Option {
	return Option{Name: RedisName, Value: cli}
}
