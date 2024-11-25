package goo_grpc

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	goo_log "github.com/liqiongtao/googo.io/goo-log"
	pb_goo_v1 "github.com/liqiongtao/googo.io/goo-proto/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// 客户端 - 单向拦截器 - 日志
func clientUnaryInterceptorLog() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		err := invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {
			goo_log.WithField("context", ctx).WithField("method", method).WithField("req", req).Error("请求失败", err)
		}
		return err
	}
}

// 客户端 - 流式拦截器 - 日志
func clientStreamInterceptorLog() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		stream, err := cc.NewStream(ctx, desc, method, opts...)
		if err != nil {
			goo_log.WithField("context", ctx).WithField("desc", desc).WithField("method", method).Error("请求失败", err)
		}
		return stream, err
	}
}

// 服务端 - 单向拦截器 - 日志
func serverUnaryInterceptorLog(noLogMethods map[string]struct{}) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		log := goo_log.WithTag("goo-grpc").WithField("method", info.FullMethod)
		ctx = context.WithValue(ctx, "log", log)

		if v, ok := req.(*pb_goo_v1.Request); ok {
			var vv interface{}
			if err = json.Unmarshal(v.Data, &vv); err == nil {
				log.WithField("request", vv)
			} else {
				log.WithField("request", req)
			}
		} else {
			log.WithField("request", req)
		}

		if md, ok := metadata.FromIncomingContext(ctx); ok {
			log.WithField("metadata", md)
		}

		var startTime = time.Now()

		defer func() {
			log.WithField("duration", fmt.Sprintf("%dms", time.Since(startTime)/1e6))

			if rst, ok := resp.(*pb_goo_v1.Response); ok && rst != nil {
				var v interface{}
				if err = json.Unmarshal(rst.Data, &v); err != nil {
					log.WithField("response", map[string]interface{}{
						"code":    rst.Code,
						"message": rst.Message,
						"data":    v,
					})
				}
			} else {
				log.WithField("response", resp)
			}

			if err != nil {
				if s, _ := status.FromError(err); s != nil {
					log.WithField("response", s.Proto())
				}
				log.Error("请求失败")
				return
			}

			// 不打印日志的方法
			if _, ok := noLogMethods[info.FullMethod]; ok {
				return
			}

			log.Debug()
		}()

		resp, err = handler(ctx, req)
		return
	}
}

// 服务端 - 流式拦截器 - 日志
func serverStreamInterceptorLog(noLogMethods map[string]struct{}) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		log := goo_log.WithTag("goo-grpc").WithField("method", info.FullMethod)

		if md, ok := metadata.FromIncomingContext(ss.Context()); ok {
			log.WithField("metadata", md)
		}

		var startTime = time.Now()

		defer func() {
			log.WithField("duration", fmt.Sprintf("%dms", time.Since(startTime)/1e6))

			if err != nil {
				log.Error("请求失败", err)
				return
			}

			// 不打印日志的方法
			if _, ok := noLogMethods[info.FullMethod]; ok {
				return
			}

			log.Debug()
		}()

		err = handler(srv, ss)
		return
	}
}
