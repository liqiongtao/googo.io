package goo

import (
	"fmt"
	"github.com/gin-gonic/gin"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
	"time"
)

// 定义控制器抽象类
type iController interface {
	DoHandle(ctx *Context) *Response
}

// 定义控制器调用实现
func Handler(controller iController) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer goo_utils.Recover()

		rsp := controller.DoHandle(&Context{Context: ctx})

		start, _ := ctx.Get("__timestamp")
		ctx.Header("X-Response-Time", fmt.Sprintf("%dms", time.Since(start.(time.Time))/1e6))

		if rsp == nil {
			return
		}

		ctx.Set("__response", rsp)
		ctx.JSON(200, rsp)
	}
}
