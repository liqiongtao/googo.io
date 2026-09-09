package goo_redis

import (
	"time"

	"github.com/go-redis/redis"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
)

type Client struct {
	Config
	*redis.Client
}

func New(conf Config) (cli *Client, err error) {
	cli = &Client{Config: conf}

	cli.Client = redis.NewClient(&redis.Options{
		Addr:     conf.Addr,
		Password: conf.Password,
		DB:       conf.DB,

		// 连接池
		PoolSize:     20,
		MinIdleConns: 5,
		PoolTimeout:  30 * time.Second,
		IdleTimeout:  5 * time.Minute,

		// 超时
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second, // 调大读超时
		WriteTimeout: 10 * time.Second,
	})

	if err = cli.Ping().Err(); err != nil {
		goo_log.WithTag("goo-redis").Error(err)
		_ = cli.Close()
		cli = nil
		return
	}

	return
}
