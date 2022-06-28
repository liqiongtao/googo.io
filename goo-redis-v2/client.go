package goo_redis

import (
	"github.com/go-redis/redis"
	goo_cron "github.com/liqiongtao/googo.io/goo-cron"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
)

type Client struct {
	*redis.Client
}

func New(conf Config) (cli *Client, err error) {
	cli = &Client{
		Client: redis.NewClient(&redis.Options{
			Addr:     conf.Addr,
			Password: conf.Password,
			DB:       conf.DB,
		}),
	}

	if err = cli.Ping().Err(); err != nil {
		goo_log.WithTag("goo-redis").Error(err)
		return
	}

	if conf.AutoPing {
		goo_cron.SecondX(5, func() {
			if err := cli.Ping().Err(); err != nil {
				goo_log.WithTag("goo-redis").Error(err)
			}
		})
	}

	return
}
