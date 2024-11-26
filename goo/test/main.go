package main

import (
	"github.com/gin-gonic/gin"
	"github.com/liqiongtao/googo.io/goo"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
)

func main() {
	s := goo.NewServer(
		goo.EnvOption(goo.DEVELOPMENT),
		goo.ServerNameOption("my-test"),
	)

	ctx := MyContext{}

	s.GET("/", goo.Handler(&ctx, MyController{}))

	s.Run("127.0.0.1:18901")
}

// 定义控制器
type MyController struct {
}

func (m MyController) DoHandle(ctx goo.Context) *goo.Response {
	ctx.Log().Debug("1")
	ctx.Log().Debug("2")
	ctx.Log().Debug("3")
	ctx.Log().Debug("4")
	return goo.Success("ok")
}

// 定义上下文
type MyContext struct {
	*gin.Context
}

func (c *MyContext) GinContext() *gin.Context {
	return c.Context
}

func (c *MyContext) SetGinContext(ctx *gin.Context) {
	c.Context = ctx
}

func (c *MyContext) Log() *goo_log.Entry {
	return goo_log.WithField("trace_id", goo.RequestId(c.Context))
}
