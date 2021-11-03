package goo_context

import (
	"context"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
	"sync"
	"syscall"
	"time"
)

var (
	cancelContext     context.Context
	cancelFunc        context.CancelFunc
	cancelContextOnce sync.Once
)

func Cancel() context.Context {
	cancelContextOnce.Do(func() {
		cancelContext, cancelFunc = context.WithCancel(context.Background())
		goo_utils.AsyncFunc(func() {
			for sig := range signalCH {
				switch sig {
				case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
					cancelFunc()
					return
				case syscall.SIGUSR1:
				case syscall.SIGUSR2:
				}
			}
		})
	})
	return cancelContext
}

func Timeout(d time.Duration, v ...map[string]interface{}) context.Context {
	var parent = context.Background()

	if l := len(v); l > 0 {
		for key, value := range v[0] {
			parent = context.WithValue(parent, key, value)
		}
	}

	ctx, cancelFunc := context.WithTimeout(parent, d)

	goo_utils.AsyncFunc(func() {
		for sig := range signalCH {
			switch sig {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				cancelFunc()
				return
			case syscall.SIGUSR1:
			case syscall.SIGUSR2:
			}
		}
	})

	return ctx
}
