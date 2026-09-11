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
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		err := invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {
			//log := goo_log.WithTag("goo-grpc").WithField("method", method).WithField("req", req)
			//if md, ok := metadata.FromIncomingContext(ctx); ok {
			//	log.WithField("metadata", md)
			//}
			//log.Error(err)
		}
		return err
	}
}

// 客户端 - 流式拦截器 - 日志
func clientStreamInterceptorLog() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		stream, err := streamer(ctx, desc, cc, method, opts...)
		if err != nil {
			//log := goo_log.WithTag("goo-grpc").WithField("method", method).WithField("desc", desc)
			//if md, ok := metadata.FromIncomingContext(ctx); ok {
			//	log.WithField("metadata", md)
			//}
			//log.Error(err)
		}
		return stream, err
	}
}

// 服务端 - 单向拦截器 - 日志
func serverUnaryInterceptorLog(noLogMethods map[string]struct{}) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		log := goo_log.WithTag("goo-grpc").WithField("method", info.FullMethod)

		if v, ok := req.(*pb_goo_v1.Request); ok {
			var vv any
			// 使用局部错误，避免污染命名返回值 err
			if uerr := json.Unmarshal(v.Data, &vv); uerr == nil {
				log = log.WithField("request", vv)
			} else {
				log = log.WithField("request", req)
			}
		} else {
			log = log.WithField("request", req)
		}

		if md, ok := metadata.FromIncomingContext(ctx); ok {
			log = log.WithField("metadata", md)
		}

		ctx = context.WithValue(ctx, "log", log)

		var startTime = time.Now()

		defer func() {
			log = log.WithField("duration", fmt.Sprintf("%dms", time.Since(startTime)/1e6))

			if rst, ok := resp.(*pb_goo_v1.Response); ok && rst != nil {
				var v any
				respData := map[string]any{
					"code":    rst.Code,
					"message": rst.Message,
				}
				if uerr := json.Unmarshal(rst.Data, &v); uerr == nil {
					respData["data"] = v
				} else {
					respData["data"] = string(rst.Data)
				}
				log = log.WithField("response", respData)
			} else if resp != nil {
				log = log.WithField("response", resp)
			}

			if err != nil {
				if s, _ := status.FromError(err); s != nil {
					log = log.WithField("response", s.Proto())
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
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		log := goo_log.WithTag("goo-grpc").WithField("method", info.FullMethod)

		if md, ok := metadata.FromIncomingContext(ss.Context()); ok {
			log = log.WithField("metadata", md)
		}

		var startTime = time.Now()

		defer func() {
			log = log.WithField("duration", fmt.Sprintf("%dms", time.Since(startTime)/1e6))

			if err != nil {
				log.Error(err)
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
