package main

import (
	"fmt"
	"github.com/liqiongtao/googo.io/goo"
	goo_http_request "github.com/liqiongtao/googo.io/goo-http-request"
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

	s.GET("/", goo.Handler(MyController{}))

	s.Run("127.0.0.1:18901")
}

// 定义控制器
type MyController struct {
}

func (m MyController) DoHandle(ctx goo.Context) *goo.Response {
	m.setName(ctx)

	ctx.Log().
		WithField("name", ctx.GetString("name")).
		WithField("request_id", ctx.GetString("request_id")).
		Debug()

	return goo.Success("ok")
}

func (m MyController) setName(ctx goo.Context) {
	ctx.Set("name", "hnatao")
	ctx.Set("request_id", ctx.Query("request_id"))
}
