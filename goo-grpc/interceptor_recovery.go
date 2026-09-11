package googrpc

import (
	"context"
	"runtime/debug"

	goolog "github.com/liqiongtao/googo.io/goo-log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// 服务端 - 单向拦截器 - panic捕获
func serverUnaryInterceptorRecovery() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		defer func() {
			if r := recover(); r != nil {
				goolog.WithTag("goo-grpc").WithField("method", info.FullMethod).
					ErrorF("panic: %v\n%s", r, debug.Stack())
				err = status.Errorf(codes.Internal, "服务异常，原因：%v", r)
			}
		}()

		resp, err = handler(ctx, req)
		return
	}
}

// 服务端 - 流式拦截器 - panic捕获
func serverStreamInterceptorRecovery() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		defer func() {
			if r := recover(); r != nil {
				goolog.WithTag("goo-grpc").WithField("method", info.FullMethod).
					ErrorF("panic: %v\n%s", r, debug.Stack())
				err = status.Errorf(codes.Internal, "服务异常，原因：%v", r)
			}
		}()

		err = handler(srv, ss)
		return
	}
}
