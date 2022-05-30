package goo

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io/ioutil"
)

// 唯一ID
func requestId(c *gin.Context) string {
	if v := c.GetHeader("X-Request-Id"); v != "" {
		return v
	}
	if v := c.Query("request_id"); v != "" {
		return v
	}
	if v := c.GetHeader("X-Trace-Id"); v != "" {
		return v
	}
	if v := c.Query("trace_id"); v != "" {
		return v
	}
	return uuid.New().String()
}

// 客户端IP
func clientIP(c *gin.Context) string {
	if v := c.GetHeader("X-Real-IP"); v != "" {
		return v
	}
	if v := c.GetHeader("X-Forwarded-For"); v != "" {
		return v
	}
	if v := c.ClientIP(); v == "::1" {
		return "127.0.0.1"
	}
	return ""
}

// 请求数据
func requestBody(c *gin.Context) interface{} {
	ctx := c.Copy()

	switch ctx.ContentType() {
	case "application/x-www-form-urlencoded", "text/xml":
		b, _ := ioutil.ReadAll(ctx.Request.Body)
		return string(b)

	case "application/json":
		b, _ := ioutil.ReadAll(ctx.Request.Body)
		var body interface{}
		if err := json.Unmarshal(b, &body); err != nil {
			return string(b)
		}
		return body
	}

	return nil
}
