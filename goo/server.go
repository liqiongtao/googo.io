package goo

import (
	"fmt"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"io/ioutil"
	"os"
	"strings"
)

type HandlerFunc func(ctx *Context)

// 定义web服务
type Server struct {
	*gin.Engine

	// 根路径
	baseDir string

	// 服务名称
	serverName string

	// 不允许访问的路径
	noAccessPath map[string]struct{}

	// 不记录日志的路径
	noLogPath map[string]struct{}

	// 跨域
	corsHeaders []string
}

func NewServer(opts ...Option) *Server {
	s := &Server{
		Engine: gin.New(),
		noAccessPath: map[string]struct{}{
			"/favicon.ico": {},
		},
		noLogPath: map[string]struct{}{
			"/favicon.ico": {},
		},
		corsHeaders: []string{
			"Content-Type", "Content-Length",
			"Accept", "Referer", "User-Agent", "Authorization",
			"X-Requested-Id", "X-Request-Timestamp", "X-Request-Sign",
			"X-Request-AppId", "X-Request-Source", "X-Request-Token",
		},
	}

	s.apply(opts...)

	s.Engine.NoRoute(s.noRoute)
	s.Engine.NoMethod(s.noMethod)

	s.Engine.Use(s.cors, s.noAccess, s.recovery)

	return s
}

// 参数设置
func (s *Server) apply(opts ...Option) {
	if l := len(opts); l == 0 {
		return
	}

	for _, opt := range opts {
		switch opt.Name {
		case baseDir:
			s.baseDir = opt.Value.(string)

		case serverName:
			s.serverName = opt.Value.(string)

		case noAccessPaths:
			for _, value := range opt.Value.([]string) {
				s.noAccessPath[value] = struct{}{}
			}

		case noLogPaths:
			for _, value := range opt.Value.([]string) {
				s.noLogPath[value] = struct{}{}
			}

		case corsHeaders:
			for _, value := range opt.Value.([]string) {
				s.corsHeaders = append(s.corsHeaders, value)
			}
		}
	}
}

// 启动服务
func (s *Server) Run(addr string) {
	pid := fmt.Sprintf("%d", os.Getpid())
	if err := ioutil.WriteFile(".pid", []byte(pid), 0644); err != nil {
		goo_log.Panic(err.Error())
	}
	endless.NewServer(addr, s.Engine).ListenAndServe()
}

// 路由
func (s *Server) Router(fn func(s *Server)) *Server {
	fn(s)
	return s
}

// 中间件
func (s *Server) Use(handlers ...HandlerFunc) *Server {
	s.Engine.Use(func(c *gin.Context) {
		ctx := NewContext(c)

		ctx.Set("__request_id", ctx.requestId())
		ctx.Set("__base_dir", s.baseDir)
		ctx.Set("__server_name", s.serverName)
		ctx.Set("__server", s)

		defer func() {
			if err := recover(); err != nil {
				ctx.AbortWithStatusJSON(200, Error(500, "请求异常"), err)
			}
		}()

		for _, handler := range handlers {
			handler(ctx)
			if ctx.IsAborted() {
				return
			}
		}
	})

	return s
}

// 跨域
func (s *Server) cors(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")
	ctx.Header("Access-Control-Allow-Headers", strings.Join(s.corsHeaders, ","))
	ctx.Next()
}

// 禁止访问
func (s *Server) noAccess(ctx *gin.Context) {
	if ctx.Request.Method == "OPTIONS" {
		ctx.AbortWithStatus(200)
		return
	}

	if _, ok := s.noAccessPath[ctx.Request.URL.Path]; ok {
		ctx.AbortWithStatus(200)
		return
	}

	ctx.Next()
}

// 捕获panic信息
func (s *Server) recovery(c *gin.Context) {
	ctx := NewContext(c)

	defer func() {
		if err := recover(); err != nil {
			rsp := Error(500, "请求异常")
			ctx.AbortWithStatusJSON(200, rsp, err)
		}
	}()

	ctx.Next()
}

// 找不到路由
func (s *Server) noRoute(c *gin.Context) {
	rsp := Error(404, "Page Not Found")
	NewContext(c).AbortWithStatusJSON(404, rsp)
}

// 找不到方法
func (*Server) noMethod(c *gin.Context) {
	rsp := Error(405, "Method not allowed")
	NewContext(c).AbortWithStatusJSON(405, rsp)
}
