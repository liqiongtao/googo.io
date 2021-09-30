package goo_grpc

import (
	"context"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
)

func New() *Server {
	return &Server{
		grpc.NewServer(grpc_middleware.WithUnaryServerChain(GRPCInterceptor)),
	}
}

type Server struct {
	*grpc.Server
}

func (s *Server) Serve(addr string) error {
	grpc_health_v1.RegisterHealthServer(s.Server, &Health{})

	return NewGRPCGraceful("tcp", addr, s.Server).Serve()
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
