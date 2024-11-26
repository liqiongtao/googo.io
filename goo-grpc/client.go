package goo_grpc

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"time"
)

func Dial(addr string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	opts = append([]grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			// 客户端在该时间内未收到任何数据时发送 ping
			// 建议设置比服务端的 Time (5分钟) 小一些，确保客户端的连接不会因为服务端的 idle 检测而断开
			Time: 4 * time.Minute,
			// ping 请求的超时时间
			// 建议设置比服务端的 Timeout (20秒) 小一些
			Timeout: 15 * time.Second,
			// 允许在没有活动流的情况下发送ping
			PermitWithoutStream: true,
		}),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(MaxRecvMsgSize), grpc.MaxCallSendMsgSize(MaxSendMsgSize)),
		grpc.WithChainUnaryInterceptor(clientUnaryInterceptorLog()),
		grpc.WithChainStreamInterceptor(clientStreamInterceptorLog()),
	}, opts...)
	return grpc.Dial(addr, opts...)
}

func DialContext(ctx context.Context, addr string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	opts = append([]grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			// 客户端在该时间内未收到任何数据时发送 ping
			// 建议设置比服务端的 Time (5分钟) 小一些，确保客户端的连接不会因为服务端的 idle 检测而断开
			Time: 4 * time.Minute,
			// ping 请求的超时时间
			// 建议设置比服务端的 Timeout (20秒) 小一些
			Timeout: 15 * time.Second,
			// 允许在没有活动流的情况下发送ping
			PermitWithoutStream: true,
		}),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(MaxRecvMsgSize), grpc.MaxCallSendMsgSize(MaxSendMsgSize)),
		grpc.WithChainUnaryInterceptor(clientUnaryInterceptorLog()),
		grpc.WithChainStreamInterceptor(clientStreamInterceptorLog()),
	}, opts...)
	return grpc.DialContext(ctx, addr, opts...)
}
