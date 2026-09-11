package googrpc

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

func dialDefaults(insecure bool, opts ...grpc.DialOption) []grpc.DialOption {
	base := []grpc.DialOption{
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                4 * time.Minute,
			Timeout:             15 * time.Second,
			PermitWithoutStream: true,
		}),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(MaxRecvMsgSize), grpc.MaxCallSendMsgSize(MaxSendMsgSize)),
		grpc.WithChainUnaryInterceptor(clientUnaryInterceptorLog()),
		grpc.WithChainStreamInterceptor(clientStreamInterceptorLog()),
	}
	if insecure {
		base = append([]grpc.DialOption{grpc.WithInsecure()}, base...)
	}
	return append(base, opts...)
}

// Dial 默认 insecure（与既有行为一致）。TLS 请用 DialSecure。
func Dial(addr string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	return grpc.Dial(addr, dialDefaults(true, opts...)...)
}

// DialSecure 不注入 insecure，由调用方提供 WithTransportCredentials。
func DialSecure(addr string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	return grpc.Dial(addr, dialDefaults(false, opts...)...)
}

func DialContext(ctx context.Context, addr string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	return grpc.DialContext(ctx, addr, dialDefaults(true, opts...)...)
}

func DialContextSecure(ctx context.Context, addr string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	return grpc.DialContext(ctx, addr, dialDefaults(false, opts...)...)
}
