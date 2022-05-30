package goo_redis

import (
	"context"
	"github.com/go-redis/redis"
	goo_context "github.com/liqiongtao/googo.io/goo-context"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"time"
)

func NewRedis(conf Config) *Redis {
	return &Redis{
		ctx:  goo_context.Cancel(),
		conf: conf,
	}
}

type Redis struct {
	ctx  context.Context
	conf Config
	*redis.Client
}

func (r *Redis) connect() (err error) {
	r.Client = redis.NewClient(&redis.Options{
		Addr:     r.conf.Addr,
		Password: r.conf.Password,
		DB:       r.conf.DB,
	})
	return r.Client.Ping().Err()
}

func (r *Redis) ping() {
	dur := 60 * time.Second
	ti := time.NewTimer(dur)
	defer ti.Stop()
	for {
		select {
		case <-r.ctx.Done():
			return
		case <-ti.C:
			if err := r.Client.Ping().Err(); err != nil {
				goo_log.WithTag("goo-redis").Error(err)
			}
			ti.Reset(dur)
		}
	}
}
