package goo

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	goo_log "googo.io/goo-log"
	goo_utils "googo.io/goo-utils"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Context struct {
	*gin.Context
}

var (
	signalCH          chan os.Signal
	signalOnce        sync.Once
	cancelContext     context.Context
	cancelFunc        context.CancelFunc
	cancelContextOnce sync.Once
)

func Signal() chan os.Signal {
	signalOnce.Do(func() {
		signalCH = make(chan os.Signal)
		signal.Notify(signalCH, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
	})
	return signalCH
}

func CancelContext() context.Context {
	cancelContextOnce.Do(func() {
		cancelContext, cancelFunc = context.WithCancel(context.Background())
		goo_utils.AsyncFunc(func() {
			for sig := range Signal() {
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

func TimeoutContext(d time.Duration, v ...map[string]interface{}) context.Context {
	var parent = context.Background()

	if l := len(v); l > 0 {
		for key, value := range v[0] {
			parent = context.WithValue(parent, key, value)
		}
	}

	ctx, cancel := context.WithTimeout(parent, d)

	goo_utils.AsyncFunc(func() {
		for sig := range Signal() {
			switch sig {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				goo_log.Error(fmt.Sprintf("request timeout in %fs", d.Seconds()))
				cancel()
				return
			case syscall.SIGUSR1:
			case syscall.SIGUSR2:
			}
		}
	})

	return ctx
}
