package goo_grpc

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthFunc func(md metadata.MD, ctx context.Context, fullMethod string) (context.Context, error)

// 服务端 - 单向拦截器 - 认证
func serverUnaryInterceptorAuth(authFunc AuthFunc) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if authFunc == nil {
			return handler(ctx, req)
		}

		md, _ := metadata.FromIncomingContext(ctx)
		ctxx, err := authFunc(md, ctx, info.FullMethod)
		if err != nil {
			return nil, status.Errorf(401, "认证失败，原因：%s", err)
		}

		return handler(ctxx, req)
	}
}

// 服务端 - 流式拦截器 - 认证
func serverStreamInterceptorAuth(authFunc AuthFunc) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if authFunc == nil {
			return handler(srv, ss)
		}

		md, _ := metadata.FromIncomingContext(ss.Context())
		ctxx, err := authFunc(md, ss.Context(), info.FullMethod)
		if err != nil {
			return status.Errorf(401, "认证失败，原因：%s", err)
		}

		ssa := newServerStreamAuth(ss)
		ssa.ctx = ctxx

		return handler(srv, ssa)
	}
}

type serverStreamAuth struct {
	grpc.ServerStream
	ctx context.Context
}

func (ssa *serverStreamAuth) Context() context.Context {
	return ssa.ctx
}

func newServerStreamAuth(ss grpc.ServerStream) *serverStreamAuth {
	if ssa, ok := ss.(*serverStreamAuth); ok {
		return ssa
	}
	return &serverStreamAuth{ServerStream: ss, ctx: ss.Context()}
}
