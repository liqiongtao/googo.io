package goo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type HandlerFunc func(ctx *Context)

type Context struct {
	*gin.Context
	*Response
}

// 定义web服务
type Server struct {
	*gin.Engine
}

func NewServer(opts ...Option) *Server {
	s := &Server{
		Engine: gin.New(),
	}

	// 不允许访问的路径
	noAccessPathMap := map[string]struct{}{
		"/favicon.ico": {},
	}

	// 不记录日志的路径
	noLogPathMap := map[string]struct{}{
		"/favicon.ico": {},
	}

	// 跨域http.header字段
	headFields := []string{
		"Content-Type", "Content-Length",
		"Accept", "Referer", "User-Agent", "Authorization",
		"X-Requested-Id", "X-Request-Timestamp", "X-Request-Sign",
		"X-Request-AppId", "X-Request-Source", "X-Request-Token",
	}

	if l := len(opts); l > 0 {
		for _, opt := range opts {
			switch opt.Name {
			case noAccessPaths:
				for _, value := range opt.Value.([]string) {
					noAccessPathMap[value] = struct{}{}
				}
			case noLogPaths:
				for _, value := range opt.Value.([]string) {
					noLogPathMap[value] = struct{}{}
				}
			case corsHeaderKeys:
				for _, value := range opt.Value.([]string) {
					headFields = append(headFields, value)
				}
			}
		}
	}

	s.Engine.NoRoute(s.noRoute())
	s.Engine.NoMethod(s.noMethod())

	s.Engine.Use(s.cors(headFields))
	s.Engine.Use(s.noAccess(noAccessPathMap))
	s.Engine.Use(s.logger(noLogPathMap))
	s.Engine.Use(s.recovery())

	return s
}

// 启动服务
// 生成.pid文件
func (s *Server) Run(addr string) {
	pid := fmt.Sprintf("%d", os.Getpid())
	if err := ioutil.WriteFile(".pid", []byte(pid), 0644); err != nil {
		goo_log.Panic(err.Error())
	}
	endless.NewServer(addr, s.Engine).ListenAndServe()
}

// 中间件
func (s *Server) Use(handlers ...HandlerFunc) {
	s.Engine.Use(func(c *gin.Context) {
		for _, handler := range handlers {
			ctx := &Context{Context: c, Response: &Response{}}
			handler(ctx)
			if code := ctx.Response.Code; code > 0 {
				ctx.Set("__response", ctx.Response)
				ctx.JSON(int(code), ctx.Response)
				ctx.Abort()
				return
			}
		}
	})
}

// 路由
func (s *Server) Router(fn func(s *Server)) {
	fn(s)
}

// 跨域
func (*Server) cors(headFields []string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")
		ctx.Header("Access-Control-Allow-Headers", strings.Join(headFields, ","))
		ctx.Next()
	}
}

// 禁止访问
func (*Server) noAccess(noAccessPathMap map[string]struct{}) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(200)
			return
		}

		if _, ok := noAccessPathMap[ctx.Request.URL.Path]; ok {
			ctx.AbortWithStatus(200)
			return
		}

		ctx.Next()
	}
}

// 日志
func (*Server) logger(noAccessLogPathMap map[string]struct{}) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		beginT := time.Now()

		requestId := (func() string {
			if v := ctx.GetHeader("X-Request-Id"); v != "" {
				return v
			} else if v := ctx.Query("request_id"); v != "" {
				return v
			} else if v := ctx.PostForm("request_id"); v != "" {
				return v
			}
			return ""
		})()

		clientIp := (func() (ip string) {
			if ip = ctx.GetHeader("X-Real-IP"); ip != "" {
				return
			}
			if ip = ctx.GetHeader("X-Forwarded-For"); ip != "" {
				return
			}
			if ip = ctx.ClientIP(); ip == "::1" {
				return "127.0.0.1"
			}
			return
		})()

		ctx.Set("__request_id", requestId)

		request := map[string]interface{}{
			"method":    ctx.Request.Method,
			"uri":       ctx.Request.RequestURI,
			"client-ip": clientIp,
		}

		defer func() {
			if _, ok := noAccessLogPathMap["*"]; ok {
				return
			}
			if _, ok := noAccessLogPathMap[ctx.Request.URL.Path]; ok {
				return
			}

			if v := ctx.GetHeader("Authorization"); v != "" {
				request["authorization"] = v
			}
			if v := ctx.GetHeader("Content-Type"); v != "" {
				request["content-type"] = v
			}
			if requestId != "" {
				request["x-request-id"] = requestId
			}

			switch ctx.ContentType() {
			case "application/x-www-form-urlencoded", "text/xml":
				buf, _ := ioutil.ReadAll(ctx.Request.Body)
				ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
				ctx.Set("__request_body", string(buf))
				request["body"] = string(buf)
			case "application/json":
				buf, _ := ioutil.ReadAll(ctx.Request.Body)
				ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
				ctx.Set("__request_body", string(buf))
				var body interface{}
				json.Unmarshal(buf, &body)
				request["body"] = body
			}

			l := goo_log.WithField("request", request).
				WithField("execute_time", fmt.Sprintf("%dms", time.Since(beginT)/1e6))

			if v, ok := ctx.Get("__response"); ok {
				l.WithField("response", v)
				if errs := v.(*Response).Errors; len(errs) > 0 {
					l.Error(errs...)
					return
				}
			}

			l.Debug()
		}()

		ctx.Next()
	}
}

// 捕获panic信息
func (*Server) recovery() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				ctx.Set("__response", Error(500, "请求异常", err))
				ctx.JSON(200, ctx.MustGet("__response"))
				ctx.Abort()
				return
			}
		}()

		ctx.Next()
	}
}

// 找不到路由
func (*Server) noRoute() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set("__response", Error(404, "Page Not Found"))
		ctx.JSON(200, ctx.MustGet("__response"))
		ctx.Abort()
	}
}

// 找不到方法
func (*Server) noMethod() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set("__response", Error(405, "Method not allowed"))
		ctx.JSON(200, ctx.MustGet("__response"))
		ctx.Abort()
	}
}
