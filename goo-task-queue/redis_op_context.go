package gootaskqueue

import (
	"context"
	"time"
)

// redis 收尾/发布等短命令超时。
const redisOpTimeout = 5 * time.Second

func redisOpContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), redisOpTimeout)
}
