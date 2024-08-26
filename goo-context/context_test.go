package goo_context

import (
	"fmt"
	"testing"
	"time"
)

func TestWithCancel(t *testing.T) {
	go func() {
		for i := 0; ; i++ {
			fmt.Println(i + 1)
			time.Sleep(time.Second)
		}
	}()

	<-WithCancel().Done()
}

func TestWithValue(t *testing.T) {
	ctx := WithLog(&Context{})
	ctx.WithValue("name", "hantao")
	ctx.WithValue("addr", "beijing")
	ctx.Log.Debug(ctx.Values())
	ctx.Log.Debug(ctx.Value("name"))
}

func TestWithLog(t *testing.T) {
	ctx := WithLog(&Context{})
	debug(ctx)
	error(ctx)

	ctx2 := WithLog(&Context{})
	debug(ctx2)
	error(ctx2)
}

func debug(ctx *Context) {
	if ctx.Log != nil {
		ctx.Log.WithTag("tag-1").WithField("name", "hnatao").Debug("this is debug")
	}
}

func error(ctx *Context) {
	if ctx.Log != nil {
		ctx.Log.WithTag("tag-2").WithField("ok", 1).Error("this is error")
	}
}
