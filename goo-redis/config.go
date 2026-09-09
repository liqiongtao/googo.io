package goo_redis

import (
	"github.com/redis/go-redis/v9"
)

type Config struct {
	Name     string `yaml:"name" json:"name"`
	Addr     string `yaml:"addr" json:"addr"`
	Password string `yaml:"password" json:"password"`
	DB       int    `yaml:"db" json:"DB"`
	Prefix   string `yaml:"prefix" json:"prefix"`
	AutoPing bool   `yaml:"auto_ping" json:"autoPing"`
	Options  *redis.Options
}

// 再导出常用类型，便于调用方无需直接依赖 go-redis 路径细节
type (
	Z         = redis.Z
	Pipeliner = redis.Pipeliner
)

var ErrNil = redis.Nil
