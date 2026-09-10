package goo_http

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

// 定义控制器抽象类
type iController interface {
	DoHandle(ctx *gin.Context) *Response
}

// 定义控制器调用实现
func Handler(controller iController) gin.HandlerFunc {
	return func(c *gin.Context) {
		beginTime := time.Now()
		resp := controller.DoHandle(c)
		opts := optsFromContext(c)

		if opts.responseHookFunc != nil {
			opts.responseHookFunc(c, resp)
		}

		if resp == nil {
			return
		}

		c.Set("__response", resp.Copy())

		// 计算执行时间
		if !beginTime.IsZero() {
			c.Header("X-Response-Duration", fmt.Sprintf("%dms", time.Since(beginTime)/1e6))
		}

		if !opts.encryptionEnable {
			c.JSON(200, resp)
			return
		}

		switch strings.ToUpper(c.Request.Method) {
		case "POST", "PUT":
		default:
			c.JSON(200, resp)
			return
		}

		if strings.Contains(strings.ToLower(c.Request.Header.Get("Content-Type")), "multipart/form-data") {
			c.JSON(200, resp)
			return
		}

		for v := range opts.encryptionExcludeUris {
			if matchURIPrefix(c.Request.RequestURI, v) {
				c.JSON(200, resp)
				return
			}
		}

		b, err := json.Marshal(&resp.Data)
		if err != nil {
			errResp := Error(5003, "数据解析失败，原因："+err.Error())
			c.Set("__response", errResp)
			c.JSON(500, errResp)
			return
		}

		enc, err := resolveEncryption(opts, c)
		if err != nil {
			errResp := Error(5004, "数据解析失败，原因："+err.Error())
			c.Set("__response", errResp)
			c.JSON(500, errResp)
			return
		}
		body, err := enc.Encode(b)
		if err != nil {
			errResp := Error(5004, "数据解析失败，原因："+err.Error())
			c.Set("__response", errResp)
			c.JSON(500, errResp)
			return
		}

		resp.Data = body
		c.JSON(200, resp)

		return
	}
}
