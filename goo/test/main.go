package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/liqiongtao/googo.io/goo"
	goo_http_request "github.com/liqiongtao/googo.io/goo-http-request"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
	"time"
)

func main() {
	goo_utils.AsyncFunc(func() {
		time.Sleep(3 * time.Second)
		for i := 0; i < 10; i++ {
			go goo_http_request.Get(fmt.Sprintf("http://127.0.0.1:18901?request_id=%d", i))
		}
	})

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
