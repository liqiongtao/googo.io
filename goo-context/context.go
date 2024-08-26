package goo_context

import (
	"context"
	"encoding/json"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"os"
	"os/signal"
	"syscall"
)

type Context struct {
	context.Context
	Log *goo_log.Entry
	v   map[string]any
}

func (ctx *Context) WithValue(key string, value any) *Context {
	if ctx.v == nil {
		ctx.v = map[string]any{}
	}
	ctx.v[key] = value
	return ctx
}

func (ctx *Context) Value(key string) any {
	if ctx.v == nil {
		ctx.v = map[string]any{}
	}
	if v, ok := ctx.v[key]; ok {
		return v
	}
	return nil
}

func (ctx *Context) Values() map[string]any {
	if ctx.v == nil {
		ctx.v = map[string]any{}
	}
	return ctx.v
}

func (ctx *Context) Json() []byte {
	v := ctx.Values()
	b, _ := json.Marshal(&v)
	return b
}

func (ctx *Context) String() string {
	return string(ctx.Json())
}

func WithCancel() *Context {
	sig := make(chan os.Signal)
	ctx, cancel := context.WithCancel(context.TODO())

	signal.Notify(sig, syscall.SIGHUP, syscall.SIGUSR1, syscall.SIGUSR2,
		syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGKILL)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				goo_log.Error(r)
			}
		}()

		for ch := range sig {
			switch ch {
			case syscall.SIGUSR1: // kill -USR1

			case syscall.SIGUSR2: // kill -USR2

			case syscall.SIGHUP: // kill -1

			case syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGKILL: // kill -9 or ctrl+c
				cancel()
			}
		}
	}()

	return &Context{Context: ctx}
}

func WithLog(ctx *Context) *Context {
	if ctx.Log == nil {
		ctx.Log = goo_log.WithTag("goo-log")
	}
	return ctx
}
