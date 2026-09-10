package goo_http

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// 唯一ID
func RequestId(c *gin.Context) string {
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
	if v := c.Query("trace-id"); v != "" {
		return v
	}
	if v := c.GetString("__trace_id"); v != "" {
		return v
	}
	traceId := uuid.New().String()
	c.Set("__trace_id", traceId)
	return traceId
}

// 客户端IP
func ClientIP(c *gin.Context) string {
	if v := c.GetHeader("X-Real-IP"); v != "" {
		return v
	}
	if v := c.GetHeader("X-Forwarded-For"); v != "" {
		return v
	}
	if v := c.ClientIP(); v == "::1" {
		return "127.0.0.1"
	} else if v != "" {
		return v
	}
	return ""
}

// 请求数据（读完后回填 Body，供后续 handler 再读）
func RequestBody(c *gin.Context) interface{} {
	contentType := c.ContentType()
	switch contentType {
	case "application/x-www-form-urlencoded", "text/xml", "application/json":
	default:
		return nil
	}

	b, err := readAndRestoreBody(c)
	if err != nil || len(b) == 0 {
		return nil
	}

	if contentType == "application/json" {
		var body interface{}
		if err := json.Unmarshal(b, &body); err == nil {
			return body
		}
	}

	return string(b)
}

// readAndRestoreBody 读取请求体并回填为可再次读取的 Reader。
func readAndRestoreBody(c *gin.Context) ([]byte, error) {
	if c.Request.Body == nil {
		return nil, nil
	}
	b, err := io.ReadAll(c.Request.Body)
	_ = c.Request.Body.Close()
	if err != nil {
		c.Request.Body = io.NopCloser(bytes.NewReader(nil))
		c.Request.ContentLength = 0
		return nil, err
	}
	c.Request.Body = io.NopCloser(bytes.NewReader(b))
	c.Request.ContentLength = int64(len(b))
	return b, nil
}

// matchURIPrefix 路径前缀匹配：/api 匹配 /api、/api/、/api?x、/api/foo，不匹配 /apiEvil；空前缀不匹配。
func matchURIPrefix(uri, prefix string) bool {
	if prefix == "" {
		return false
	}
	if uri == prefix {
		return true
	}
	if !strings.HasPrefix(uri, prefix) {
		return false
	}
	if strings.HasSuffix(prefix, "/") {
		return true
	}
	switch uri[len(prefix)] {
	case '/', '?':
		return true
	default:
		return false
	}
}
