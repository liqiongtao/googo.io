package goo_grpc

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"time"
)

func Dial(addr string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	opts = append(opts,
		grpc.WithInsecure(),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                10 * time.Second,
			Timeout:             100 * time.Millisecond,
			PermitWithoutStream: true,
		}),
		grpc.WithChainUnaryInterceptor(clientUnaryInterceptorLog()),
		grpc.WithChainStreamInterceptor(clientStreamInterceptorLog()),
	)
	return grpc.Dial(addr, opts...)
}

func DialContext(ctx context.Context, addr string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	opts = append(opts,
		grpc.WithInsecure(),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                10 * time.Second,
			Timeout:             100 * time.Millisecond,
			PermitWithoutStream: true,
		}),
		grpc.WithChainUnaryInterceptor(clientUnaryInterceptorLog()),
		grpc.WithChainStreamInterceptor(clientStreamInterceptorLog()),
	)
	return grpc.DialContext(ctx, addr, opts...)
}
