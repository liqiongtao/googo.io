package goo

import (
	"github.com/gin-gonic/gin"
)

// 定义控制器抽象类
type iController interface {
	DoHandle(ctx *Context) *Response
}

// 定义控制器调用实现
func Handler(controller iController) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := NewContext(c)
		rsp := controller.DoHandle(ctx)

		if rsp == nil {
			return
		}

		ctx.JSON(200, rsp, rsp.Errors...)
	}
}
