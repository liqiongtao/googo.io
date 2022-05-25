package goo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"io/ioutil"
	"strconv"
	"sync"
	"time"
)

type Context struct {
	*gin.Context
	s         *Server
	request   map[string]interface{}
	beginTime time.Time
}

var muNewContext sync.Mutex

func NewContext(c *gin.Context) (ctx *Context) {
	muNewContext.Lock()
	defer muNewContext.Unlock()

	if v, ok := c.Get("__context"); ok {
		ctx = v.(*Context)
		return
	}

	ctx = &Context{
		Context:   c,
		beginTime: time.Now(),
	}
	ctx.request = ctx.requestInfo()

	ctx.Context.Set("__context", ctx)

	return
}

// 返回json数据
func (ctx *Context) JSON(code int, rsp *Response, v ...interface{}) {
	executeTime := fmt.Sprintf("%dms", time.Since(ctx.beginTime)/1e6)

	l := goo_log.WithField("execute_time", executeTime).
		WithField("request", ctx.request).
		WithField("response", rsp)
	if rsp.Code > 0 {
		l.Error(v...)
	} else {
		l.Debug()
	}

	ctx.Context.Header("X-Response-Time", executeTime)
	ctx.Context.JSON(code, rsp)
}

// 中止并返回json数据
func (ctx *Context) AbortWithStatusJSON(code int, rsp *Response, v ...interface{}) {
	executeTime := fmt.Sprintf("%dms", time.Since(ctx.beginTime)/1e6)

	goo_log.WithField("request", ctx.request).
		WithField("response", rsp).
		WithField("execute_time", executeTime).
		Error(v...)

	ctx.Context.Header("X-Response-Time", executeTime)
	ctx.Context.AbortWithStatusJSON(code, rsp)
}

// 请求信息
func (ctx *Context) requestInfo() (data map[string]interface{}) {
	data = map[string]interface{}{
		"method": ctx.Context.Request.Method,
		"uri":    ctx.Context.Request.RequestURI,
	}

	if v := ctx.Context.GetHeader("Authorization"); v != "" {
		data["authorization"] = v
	}
	if v := ctx.Context.GetHeader("Content-Type"); v != "" {
		data["content-type"] = v
	}
	if v := ctx.clientIP(); v != "" {
		data["client-ip"] = v
	}
	if v := ctx.requestId(); v != "" {
		data["request-id"] = v
	}
	if v := ctx.requestBody(); v != nil {
		data["body"] = v
	}

	return
}

// 唯一ID
func (ctx *Context) requestId() string {
	if v := ctx.Context.GetHeader("X-Request-Id"); v != "" {
		return v
	}
	if v := ctx.Context.Query("request_id"); v != "" {
		return v
	}
	if v := ctx.Context.PostForm("request_id"); v != "" {
		return v
	}
	nw := time.Now()
	return nw.Format("20060102-150304-") + strconv.Itoa(nw.Nanosecond())
}

// 客户端IP
func (ctx *Context) clientIP() string {
	if v := ctx.Context.GetHeader("X-Real-IP"); v != "" {
		return v
	}
	if v := ctx.Context.GetHeader("X-Forwarded-For"); v != "" {
		return v
	}
	if v := ctx.Context.ClientIP(); v == "::1" {
		return "127.0.0.1"
	}
	return ""
}

// 请求数据
// 因为数据流读取一次就关闭了，所以要放在请求参数解析之前调用
func (ctx *Context) requestBody() interface{} {
	switch ctx.ContentType() {
	case "application/x-www-form-urlencoded", "text/xml":
		buf, _ := ioutil.ReadAll(ctx.Context.Request.Body)
		ctx.Context.Request.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
		return string(buf)

	case "application/json":
		buf, _ := ioutil.ReadAll(ctx.Context.Request.Body)
		ctx.Context.Request.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
		var body interface{}
		if err := json.Unmarshal(buf, &body); err != nil {
			return string(buf)
		}
		return body
	}

	return nil
}
