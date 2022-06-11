package goo_grpc

import (
	"context"
	"encoding/json"
	"fmt"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	pb_goo_v1 "github.com/liqiongtao/googo.io/goo-proto/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"time"
)

// 客户端 - 单向拦截器 - 日志
func clientUnaryInterceptorLog() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// 客户端 - 流式拦截器 - 日志
func clientStreamInterceptorLog() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		return cc.NewStream(ctx, desc, method, opts...)
	}
}

// 服务端 - 单向拦截器 - 日志
func serverUnaryInterceptorLog() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		l := goo_log.WithTag("goo-grpc").
			WithField("method", info.FullMethod)

		if v, ok := req.(*pb_goo_v1.Request); ok {
			var vv interface{}
			if err = json.Unmarshal(v.Data, &vv); err == nil {
				l.WithField("request", vv)
			} else {
				l.WithField("request", req)
			}
		} else {
			l.WithField("request", req)
		}

		if md, ok := metadata.FromIncomingContext(ctx); ok {
			l.WithField("metadata", md)
		}

		var startTime = time.Now()

		defer func() {
			l.WithField("duration", fmt.Sprintf("%dms", time.Since(startTime)/1e6))

			if rst, ok := resp.(*pb_goo_v1.Response); ok && rst != nil {
				var v interface{}
				if err := json.Unmarshal(rst.Data, &v); err != nil {
					l.WithField("response", map[string]interface{}{
						"code":    rst.Code,
						"message": rst.Message,
						"data":    v,
					})
				}
			} else {
				l.WithField("response", resp)
			}

			if err == nil {
				l.Debug()
				return
			}

			if s, _ := status.FromError(err); s != nil {
				l.WithField("response", s.Proto())
			}

			l.Error()
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
			l.WithField("duration", fmt.Sprintf("%dms", time.Since(startTime)/1e6))

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
