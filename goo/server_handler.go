package goo

import (
	"github.com/gin-gonic/gin"
)

// 定义控制器抽象类
type iController interface {
	DoHandle(c *gin.Context) *Response
}

// 定义控制器调用实现
func Handler(controller iController) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp := controller.DoHandle(c)

		if resp == nil {
			return
		}

		c.Set("__response", resp)

		c.JSON(200, resp)
	}
}
