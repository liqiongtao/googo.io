package goo

import (
	"bytes"
	"fmt"
	"github.com/fvbock/endless"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

// 定义web服务
type Server struct {
	*gin.Engine
}

func NewServer(opt ...Option) *Server {
	for _, o := range opt {
		o.apply(defaultOptions)
	}

	s := &Server{
		Engine: gin.New(),
	}

	s.Engine.NoRoute(s.noRoute)
	s.Engine.NoMethod(s.noMethod)

	s.Use(s.cors, s.noAccess, s.setFields, s.encrypt, s.log, s.recovery)

	return s
}

// 启动服务
func (s *Server) Run(addr string) {
	pid := fmt.Sprintf("%d", os.Getpid())
	if err := ioutil.WriteFile(".pid", []byte(pid), 0644); err != nil {
		goo_log.Panic(err.Error())
	}

	// 性能分析
	if defaultOptions.pprofEnable {
		pprof.Register(s.Engine, "/goo/pprof")
	}

	endless.NewServer(addr, s.Engine).ListenAndServe()
}

// 跨域
func (s *Server) cors(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")
	c.Header("Access-Control-Allow-Headers", strings.Join(defaultOptions.corsHeaders, ","))
	c.Next()
}

// 禁止访问
func (s *Server) noAccess(c *gin.Context) {
	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(200)
		return
	}

	if _, ok := defaultOptions.noAccessPath[c.Request.URL.Path]; ok {
		c.AbortWithStatus(200)
		return
	}

	c.Next()
}

// 设置字段
func (s *Server) setFields(c *gin.Context) {
	c.Set("__begin_time", time.Now())
	c.Set("__server_name", defaultOptions.serverName)
	c.Set("__env", defaultOptions.env.Tag())

	c.Next()
}

// 加解密
func (s *Server) encrypt(c *gin.Context) {
	if !defaultOptions.encryptionEnable {
		c.Next()
		return
	}

	switch strings.ToUpper(c.Request.Method) {
	case "POST", "PUT":
	default:
		c.Next()
		return
	}

	switch strings.ToLower(c.Request.Header.Get("Content-Type")) {
	case "multipart/form-data":
		c.Next()
		return
	}

	if _, ok := defaultOptions.encryptionExcludeUris[c.Request.RequestURI]; ok {
		c.Next()
		return
	}

	var buf bytes.Buffer
	io.Copy(&buf, c.Request.Body)

	b, err := defaultOptions.encryption.Decode(buf.String())
	if err != nil {
		s.abortWithStatus50X(c, 5002, "解码失败，原因："+err.Error())
		return
	}

	if l := len(b); l == 0 {
		c.Next()
		return
	}

	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(b))

	c.Next()
}

// log
func (s *Server) log(c *gin.Context) {
	if _, ok := defaultOptions.noLogPath[c.Request.RequestURI]; ok {
		c.Next()
		return
	}

	beginTime := c.GetTime("__begin_time")

	header := gin.H{}
	if v := c.GetHeader("Authorization"); v != "" {
		header["authorization"] = v
	}
	if v := c.GetHeader("Content-Type"); v != "" {
		header["content-type"] = v
	}

	req := gin.H{
		"method":    c.Request.Method,
		"uri":       c.Request.RequestURI,
		"header":    header,
		"client-ip": clientIP(c),
		"trace-id":  requestId(c),
	}
	if v := requestBody(c); v != nil {
		req["body"] = v
	}

	l := goo_log.WithTag("goo-api").
		WithField("request", req)

	c.Next()

	if !beginTime.IsZero() {
		l.WithField("duration", fmt.Sprintf("%dms", time.Since(beginTime)/1e6))
	}

	if resp, has := c.Get("__response"); has {
		l.WithField("response", resp)
		if r, ok := resp.(*Response); ok {
			if ll := len(r.Errors); ll > 0 {
				l.Error(r.Errors)
				return
			}

			if r.Code > 0 {
				l.Warn()
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
