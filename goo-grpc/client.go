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
			Time:                5 * time.Minute,
			Timeout:             20 * time.Second,
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
			Time:                5 * time.Minute,
			Timeout:             20 * time.Second,
			PermitWithoutStream: true,
		}),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(MaxRecvMsgSize), grpc.MaxCallSendMsgSize(MaxSendMsgSize)),
		grpc.WithChainUnaryInterceptor(clientUnaryInterceptorLog()),
		grpc.WithChainStreamInterceptor(clientStreamInterceptorLog()),
	}, opts...)
	return grpc.DialContext(ctx, addr, opts...)
}
