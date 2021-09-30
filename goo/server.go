package goo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	goo_log "googo.io/goo-log"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

// 定义web服务
type server struct {
	*gin.Engine
}

func NewServer(opts ...Option) *server {
	s := &server{
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

	s.Use(s.ts())
	s.Use(s.cors(headFields))
	s.Use(s.noAccess(noAccessPathMap))
	s.Use(s.logger(noLogPathMap))
	s.Use(s.recovery())

	s.NoRoute(s.noRoute())

	return s
}

// 启动服务
// 生成.pid文件
func (s *server) Run(addr string) {
	pid := fmt.Sprintf("%d", os.Getpid())
	if err := ioutil.WriteFile(".pid", []byte(pid), 0644); err != nil {
		goo_log.Panic(err.Error())
	}
	endless.NewServer(addr, s.Engine).ListenAndServe()
}

// 执行时间
func (*server) ts() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set("__timestamp", time.Now())
	}
}

// 跨域
func (*server) cors(headFields []string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")
		ctx.Header("Access-Control-Allow-Headers", strings.Join(headFields, ","))
		ctx.Next()
	}
}

// 禁止访问
func (*server) noAccess(noAccessPathMap map[string]struct{}) gin.HandlerFunc {
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
func (*server) logger(noAccessLogPathMap map[string]struct{}) gin.HandlerFunc {
	var (
		requestId = 0
		body      interface{}
	)

	return func(ctx *gin.Context) {
		if _, ok := noAccessLogPathMap[ctx.Request.URL.Path]; ok {
			ctx.Next()
			return
		}

		requestId++
		ctx.Set("__request_id", requestId)

		switch ctx.ContentType() {
		case "application/x-www-form-urlencoded", "text/xml":
			buf, _ := ioutil.ReadAll(ctx.Request.Body)
			ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
			ctx.Set("__request_body", string(buf))
			body = string(buf)
		case "application/json":
			buf, _ := ioutil.ReadAll(ctx.Request.Body)
			ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
			ctx.Set("__request_body", string(buf))
			json.Unmarshal(buf, &body)
		default:
			body = ""
		}

		ctx.Next()

		defer func() {
			start, _ := ctx.Get("__timestamp")

			clientIp := func(ctx *gin.Context) string {
				if ip := ctx.GetHeader("X-Real-IP"); ip != "" {
					return ip
				}
				if ip := ctx.GetHeader("X-Forwarded-For"); ip != "" {
					return ip
				}
				ip := ctx.ClientIP()
				if ip == "::1" {
					return "127.0.0.1"
				}
				return ip
			}(ctx)

			l := goo_log.WithField("request-id", requestId).
				WithField("request-method", ctx.Request.Method).
				WithField("request-uri", ctx.Request.RequestURI).
				WithField("request-body", body).
				WithField("client-ip", clientIp).
				WithField("authorization", ctx.GetHeader("Authorization")).
				WithField("content-type", ctx.GetHeader("Content-Type")).
				WithField("x-request-id", ctx.GetHeader("X-Request-Id")).
				WithField("x-request-sign", ctx.GetHeader("X-Request-Sign")).
				WithField("request-time", fmt.Sprintf("%dms", time.Since(start.(time.Time))/1e6))

			if rsp, _ := ctx.Get("__response"); rsp != nil {
				l.WithField("response", rsp)
				if errMsg := rsp.(*Response).ErrMsg; len(errMsg) > 0 {
					l.Error(errMsg...)
					return
				}
			}

			l.Debug()
		}()
	}
}

// 捕获panic信息
func (*server) recovery() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				rsp := Error(500, "请求异常", err)
				ctx.Set("__response", rsp)
				ctx.AbortWithStatusJSON(200, rsp)
			}
		}()

		ctx.Next()
	}
}

// 找不到路由
func (*server) noRoute() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.AbortWithStatusJSON(200, Error(404, "Page Not Found"))
	}
}
