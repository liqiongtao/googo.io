package goo_redis

import (
	"context"
	"github.com/go-redis/redis"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"time"
)

func NewRedis(ctx context.Context, config Config) *Redis {
	return &Redis{
		ctx:    ctx,
		config: config,
	}
}

type Redis struct {
	ctx    context.Context
	config Config
	*redis.Client
}

func (r *Redis) connect() (err error) {
	r.Client = redis.NewClient(&redis.Options{
		Addr:     r.config.Addr,
		Password: r.config.Password,
		DB:       r.config.DB,
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
