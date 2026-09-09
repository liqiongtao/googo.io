package goo_http

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/cloudflare/tableflip"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	goo_pprof "github.com/liqiongtao/googo.io/goo-pprof"
	"github.com/liqiongtao/googo.io/goocontext"
)

// 定义web服务
type Server struct {
	*gin.Engine
	opts      *options
	hooksOnce sync.Once
}

func NewServer(opt ...Option) *Server {
	opts := newDefaultOptions()
	for _, o := range opt {
		o.apply(opts)
	}

	s := &Server{
		Engine: gin.New(),
		opts:   opts,
	}

	s.Engine.NoRoute(s.noRoute)
	s.Engine.NoMethod(s.noMethod)

	s.Use(s.injectOpts, s.cors, s.noAccess, s.setFields, s.recovery, s.encrypt, s.log)

	return s
}

// 启动服务
func (s *Server) Run(addr string) {
	pid := fmt.Sprintf("%d", os.Getpid())
	if err := os.WriteFile(".pid", []byte(pid), 0644); err != nil {
		goo_log.Panic(err.Error())
	}

	upg, err := tableflip.New(tableflip.Options{})
	if err != nil {
		goo_log.Panic(err.Error())
	}
	defer upg.Stop()

	httpServer := &http.Server{
		Handler:           s.Engine,
		ReadHeaderTimeout: 10 * time.Second,
	}

	// 先注册钩子再 Listen/Root，避免信号窗口内空钩子直接 cancel
	s.hooksOnce.Do(func() {
		goocontext.OnRestart(func() {
			if err := upg.Upgrade(); err != nil {
				goo_log.Error(err.Error())
				return
			}
			goo_log.Warn("服务重启")
		})
		goocontext.OnExit(func() {
			goo_pprof.StopDefault()
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			if err := httpServer.Shutdown(ctx); err != nil {
				goo_log.Error(err.Error())
			}
			goo_log.Warn("服务退出")
		})
		goo_pprof.RegisterSignal()
	})

	// 性能分析：HTTP endpoint + 进程级 pprof（USR1 统一由 goo_pprof 注册）
	if s.opts.pprofEnable {
		pprof.Register(s.Engine, "/goo/pprof")
		goo_pprof.StartDefault()
	}

	lis, err := upg.Listen("tcp", addr)
	if err != nil {
		goo_log.Panic(err.Error())
	}

	go func() {
		if err := httpServer.Serve(lis); err != nil && err != http.ErrServerClosed {
			goo_log.Error(err.Error())
		}
	}()

	goo_log.InfoF("server running, addr=%s pid=%s", lis.Addr().String(), pid)

	// 通知父进程：本进程已就绪可接管
	if err := upg.Ready(); err != nil {
		goo_log.Error(err.Error())
	}

	// 父进程在子进程 Ready 后会收到 Exit，走统一优雅退出路径
	go func() {
		<-upg.Exit()
		_ = syscall.Kill(os.Getpid(), syscall.SIGTERM)
	}()

	<-goocontext.Root().Done()
}

const gooOptsKey = "__goo_opts"

func (s *Server) injectOpts(c *gin.Context) {
	c.Set(gooOptsKey, s.opts)
	c.Next()
}

func optsFromContext(c *gin.Context) *options {
	if v, ok := c.Get(gooOptsKey); ok {
		if opts, ok := v.(*options); ok && opts != nil {
			return opts
		}
	}
	return newDefaultOptions()
}

// 跨域
func (s *Server) cors(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")
	c.Header("Access-Control-Allow-Headers", strings.Join(s.opts.corsHeaders, ","))
	c.Next()
}

// 禁止访问
func (s *Server) noAccess(c *gin.Context) {
	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(200)
		return
	}

	if _, ok := s.opts.noAccessPath[c.Request.URL.Path]; ok {
		c.AbortWithStatus(200)
		return
	}

	c.Next()
}

// 设置字段
func (s *Server) setFields(c *gin.Context) {
	c.Next()
}

// 加解密
func (s *Server) encrypt(c *gin.Context) {
	if !s.opts.encryptionEnable {
		c.Next()
		return
	}

	switch strings.ToUpper(c.Request.Method) {
	case "POST", "PUT":
	default:
		c.Next()
		return
	}

	if strings.Contains(strings.ToLower(c.Request.Header.Get("Content-Type")), "multipart/form-data") {
		c.Next()
		return
	}

	for v := range s.opts.encryptionExcludeUris {
		if v == c.Request.RequestURI || strings.HasPrefix(c.Request.RequestURI, v) {
			c.Next()
			return
		}
	}

	raw, err := io.ReadAll(c.Request.Body)
	_ = c.Request.Body.Close()
	if err != nil {
		s.abortWithStatus50X(c, 5002, "读取请求失败，原因："+err.Error())
		return
	}

	b, err := s.opts.encryptionFn(c).Decode(string(raw))
	if err != nil {
		s.abortWithStatus50X(c, 5002, "解码失败，原因："+err.Error())
		return
	}

	c.Request.Body = io.NopCloser(bytes.NewReader(b))
	c.Request.ContentLength = int64(len(b))

	c.Next()
}

// log
func (s *Server) log(c *gin.Context) {
	if _, ok := s.opts.noLogPath[c.Request.URL.Path]; ok {
		c.Next()
		return
	}
	if _, ok := s.opts.noLogPath[c.Request.RequestURI]; ok {
		c.Next()
		return
	}

	beginTime := time.Now()

	header := gin.H{}
	if v := c.GetHeader("Authorization"); v != "" {
		header["authorization"] = v
	}
	if v := c.GetHeader("Content-Type"); v != "" {
		header["content-type"] = v
	}

	req := gin.H{
		"method": c.Request.Method,
		"uri":    c.Request.RequestURI,
		"header": header,
	}
	if v := RequestBody(c); v != nil {
		req["body"] = v
	}

	l := goo_log.WithTag("goo-api").
		WithField("client-ip", ClientIP(c)).
		WithField("trace-id", RequestId(c)).
		WithField("request", req)

	c.Next()

	if !beginTime.IsZero() {
		l = l.WithField("duration", fmt.Sprintf("%dms", time.Since(beginTime)/1e6))
	}

	ctx := c.Copy()
	for k, v := range ctx.Keys {
		if strings.HasPrefix(k, "__") {
			continue
		}
		l = l.WithField(k, v)
	}

	if resp, has := ctx.Get("__response"); has {
		if r, ok := resp.(*Response); resp != nil && ok {
			l = l.WithField("response", resp)
			if r == nil {
				l.Error(resp)
				return
			}

			if r.Errors != nil && len(r.Errors) > 0 {
				l.Error(r.Errors)
				return
			}

			if r.Code > 0 {
				l.Error()
				return
			}
		}
	}

	l.Debug()
}

// 捕获panic信息
func (s *Server) recovery(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			s.abortWithStatus50X(c, 5001, fmt.Sprintf("请求异常, 提示信息: %v", r))
		}
	}()

	c.Next()
}

// 找不到路由
func (s *Server) noRoute(c *gin.Context) {
	s.abortWithStatus40X(c, 404, "Page Not Found")
}

// 找不到方法
func (s *Server) noMethod(c *gin.Context) {
	s.abortWithStatus40X(c, 405, "Method not allowed")
}

func (*Server) abortWithStatus40X(c *gin.Context, code int32, msg string) {
	resp := Error(code, msg, msg)
	c.Set("__response", resp)
	c.AbortWithStatusJSON(int(code), resp)
}

func (*Server) abortWithStatus50X(c *gin.Context, code int32, msg string) {
	resp := Error(code, msg, msg)
	c.Set("__response", resp)
	c.AbortWithStatusJSON(500, resp)
}
