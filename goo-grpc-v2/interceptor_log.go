package main

import (
	"context"
	"fmt"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"time"
)

// 服务端 - 单向拦截器 - 日志
func serverUnaryInterceptorLog() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		l := goo_log.WithTag("goo-grpc").
			WithField("method", info.FullMethod).
			WithField("request", req)

		if md, ok := metadata.FromIncomingContext(ctx); ok {
			l.WithField("metadata", md)
		}

		var startTime = time.Now()

		defer func() {
			l.WithField("response", resp)
			l.WithField("execute_time", fmt.Sprintf("%dms", time.Since(startTime)/1e6))

			if err != nil {
				l.Error(err)
				return
			}
			l.Debug()
		}()

		resp, err = handler(ctx, req)
		return
	}
}

// 服务端 - 流式拦截器 - 日志
func serverStreamInterceptorLog() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		l := goo_log.WithTag("goo-grpc").
			WithField("method", info.FullMethod)

		if md, ok := metadata.FromIncomingContext(ss.Context()); ok {
			l.WithField("metadata", md)
		}

		var startTime = time.Now()

		defer func() {
			l.WithField("execute_time", fmt.Sprintf("%dms", time.Since(startTime)/1e6))

			if err != nil {
				l.Error(err)
				return
			}
			l.Debug()
		}()

		err = handler(srv, ss)
		return
	}
}
