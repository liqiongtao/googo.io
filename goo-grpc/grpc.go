package goo_grpc

import (
	"context"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	goo_etcd "github.com/liqiongtao/googo.io/goo-etcd"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
)

func New(cfg Config) *Server {
	s := &Server{
		cfg:    cfg,
		Server: grpc.NewServer(grpc_middleware.WithUnaryServerChain(GRPCInterceptor)),
	}
	return s
}

type Server struct {
	*grpc.Server
	cfg Config
}

func (s *Server) Register2Etcd(cli *goo_etcd.Client) *Server {
	var key string
	if str := s.cfg.ENV; str != "" {
		key += "/" + str
	}
	if str := s.cfg.ServiceName; str != "" {
		key += "/" + str
	}
	if str := s.cfg.Version; str != "" {
		key += "/" + str
	}
	if key == "" {
		key = s.cfg.Addr
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
	lg := goo_log.WithField("grpc-method", info.FullMethod).WithField("grpc-request", req)
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		for key, val := range md {
			lg.WithField(key, val)
		}
	}
	defer func() {
		if e := recover(); e != nil {
			lg.Error(fmt.Sprintf("%v", e))
		}
	}()
	rsp, err = handler(ctx, req)
	lg.WithField("grpc-response", rsp)
	if err == nil {
		lg.Info()
		return
	}
	lg.Error(err.Error())
	return
}
