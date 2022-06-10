package goo

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
	"time"
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

		// 计算执行时间
		beginTime := c.GetTime("__begin_time")
		if !beginTime.IsZero() {
			c.Header("duration", fmt.Sprintf("%dms", time.Since(beginTime)/1e6))
		}

		// 如果没有启用加密传输
		if !defaultOptions.enableEncryption {
			c.JSON(200, resp)
			return
		}

		// 启用加密传输
		data := gin.H{
			"code":    resp.Code,
			"message": resp.Message,
			"ts":      time.Now().Unix(),
		}
		if resp.Data != nil {
			b, _ := json.Marshal(&resp.Data)
			data["data"] = goo_utils.Base59Encoding(b, defaultOptions.encryptionKey)
		}
		c.JSON(200, data)
	}
}
