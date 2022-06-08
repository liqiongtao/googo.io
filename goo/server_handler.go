package goo

import (
	"fmt"
	"github.com/gin-gonic/gin"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
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

		if defaultOptions.enableEncodeResponse {
			if resp.Data != nil {
				resp.Data = goo_utils.Base59Encoding(fmt.Sprintf("%v", resp.Data), defaultOptions.encodeKey)
			}
		}

		c.JSON(200, resp)
	}
}
