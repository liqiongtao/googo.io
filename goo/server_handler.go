package goo

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

// 定义控制器抽象类
type iController interface {
	DoHandle(ctx *Context) *Response
}

// 定义控制器调用实现
func Handler(controller iController) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		beginT := time.Now()
		
		rsp := controller.DoHandle(&Context{Context: ctx})

		if rsp == nil {
			return
		}

		ctx.Header("X-Response-Time", fmt.Sprintf("%dms", time.Since(beginT)/1e6))

		ctx.Set("__response", rsp)
		ctx.JSON(200, rsp)
	}
}
