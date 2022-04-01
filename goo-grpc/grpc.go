package goo_grpc

import (
	"context"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	goo_etcd "github.com/liqiongtao/googo.io/goo-etcd"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"strings"
	"time"
)

func New(cfg Config) *Server {
	// 性能分析
	if cfg.DebugAddr != "" {
		goo_utils.AsyncFunc(func() {
			r := gin.Default()
			pprof.Register(r, "/goo-grpc/pprof")
			r.Run(cfg.DebugAddr)
		})
	}

	s := &Server{
		cfg: cfg,
		Server: grpc.NewServer(
			grpc_middleware.WithUnaryServerChain(GRPCInterceptor),
			grpc.KeepaliveParams(keepalive.ServerParameters{MaxConnectionIdle: 5 * time.Minute}),
		),
	}
	return s
}

type Server struct {
	*grpc.Server
	cfg Config
}

func (s *Server) Register2Etcd(cli *goo_etcd.Client) *Server {
	key := s.cfg.ServiceName
	if str := s.cfg.ENV; str != "" {
		key += "/" + str
	}
	if str := s.cfg.Version; str != "" {
		key += "/" + str
	}
	if key == "" {
		key = s.cfg.Addr
	}
	if strings.Index(key, "/") != 0 {
		key += "/" + key
	}
	cli.SetWithKeepAlive(key, s.cfg.Addr, 15)
	return s
}

func (s *Server) Serve() error {
	grpc_health_v1.RegisterHealthServer(s.Server, &Health{})
	return NewGRPCGraceful("tcp", s.cfg.Addr, s.Server).Serve()
}

func GRPCInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (rsp interface{}, err error) {
	if info.FullMethod == "/grpc.health.v1.Health/Check" {
		return
	}

	l := goo_log.WithTag("goo-grpc").
		WithField("method", info.FullMethod).
		WithField("request", req)

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		for key, val := range md {
			l.WithField(key, val)
		}
	}

	defer func() {
		if rsp != nil {
			l.WithField("response", rsp)
		}

		if e := recover(); e != nil {
			l.WithTrace().ErrorF("%v", e)
			return
		}

		if err != nil {
			l.WithTrace().Error(err)
			return
		}

		l.Debug()
	}()

	rsp, err = handler(ctx, req)

	return
}
