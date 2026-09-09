package goo_http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ErrBodyTooLarge 请求体超过 maxBodyBytes
var ErrBodyTooLarge = errors.New("request body too large")

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

func maxBodyBytesFromContext(c *gin.Context) int64 {
	opts := optsFromContext(c)
	if opts != nil && opts.maxBodyBytes > 0 {
		return opts.maxBodyBytes
	}
	return 32 << 20
}

// readBodyLimited 读取 body，超过 limit 返回错误
func readBodyLimited(r io.Reader, limit int64) ([]byte, error) {
	if limit <= 0 {
		limit = 32 << 20
	}
	b, err := io.ReadAll(io.LimitReader(r, limit+1))
	if err != nil {
		return nil, err
	}
	if int64(len(b)) > limit {
		return nil, fmt.Errorf("%w, max %d bytes", ErrBodyTooLarge, limit)
	}
	return b, nil
}

// readAndRestoreBody 读取请求体并回填为可再次读取的 Reader。
func readAndRestoreBody(c *gin.Context) ([]byte, error) {
	if c.Request.Body == nil {
		return nil, nil
	}
	b, err := readBodyLimited(c.Request.Body, maxBodyBytesFromContext(c))
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
