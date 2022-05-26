package main

import (
	"context"
	"google.golang.org/grpc"
)

type AuthFunc func(ctx context.Context, fullMethod string) (context.Context, error)

// 服务端 - 单向拦截器 - 认证
func unaryServerInterceptorAuth(authFunc AuthFunc) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if authFunc == nil {
			return handler(ctx, req)
		}

		ctxx, err := authFunc(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}

		return handler(ctxx, req)
	}
}

// 服务端 - 流式拦截器 - 认证
func streamServerInterceptorAuth(authFunc AuthFunc) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if authFunc == nil {
			return handler(srv, ss)
		}

		ctxx, err := authFunc(ss.Context(), info.FullMethod)
		if err != nil {
			return err
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
