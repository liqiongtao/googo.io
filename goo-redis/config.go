package goo_redis

import (
	"time"

	"github.com/go-redis/redis"
)

var (
	DefaultOptions = &redis.Options{
		// 连接池
		PoolSize:     20,
		MinIdleConns: 5,
		PoolTimeout:  30 * time.Second,
		IdleTimeout:  5 * time.Minute,

		// 超时
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second, // 调大读超时
		WriteTimeout: 10 * time.Second,
	}
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
