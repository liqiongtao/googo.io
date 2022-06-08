package goo

import (
	"encoding/json"
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

		if !defaultOptions.enableEncryption {
			c.JSON(200, resp)
			return
		}

		data := gin.H{
			"code":    resp.Code,
			"message": resp.Message,
		}
		if resp.Data != nil {
			b, _ := json.Marshal(&resp.Data)
			data["data"] = goo_utils.Base59Encoding(string(b), defaultOptions.encryptionKey)
		}
		c.JSON(200, data)
	}
}
