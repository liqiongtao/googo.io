package goo_grpc

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// 服务端 - 单向拦截器 - panic捕获
func serverUnaryInterceptorRecovery() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				err = status.Errorf(500, "服务异常，原因：%v", r)
			}
		}()

		resp, err = handler(ctx, req)
		return
	}
}

// 服务端 - 流式拦截器 - panic捕获
func serverStreamInterceptorRecovery() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = status.Errorf(500, "服务异常，原因：%v", r)
			}
		}()

		err = handler(srv, ss)
		return
	}
}
