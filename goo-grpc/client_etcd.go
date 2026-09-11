package googrpc

import (
	"context"
	"fmt"
	"time"

	gooetcd "github.com/liqiongtao/googo.io/goo-etcd"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

func DialWithEtcd(serviceName string, cli *gooetcd.Client, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	return dialWithEtcd(context.Background(), serviceName, cli, true, false, opts...)
}

// DialSecureWithEtcd 不注入 insecure，由调用方提供 WithTransportCredentials。
func DialSecureWithEtcd(serviceName string, cli *gooetcd.Client, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	return dialWithEtcd(context.Background(), serviceName, cli, false, false, opts...)
}

func DialContextWithEtcd(ctx context.Context, serviceName string, cli *gooetcd.Client, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	return dialWithEtcd(ctx, serviceName, cli, true, true, opts...)
}

func DialContextSecureWithEtcd(ctx context.Context, serviceName string, cli *gooetcd.Client, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	return dialWithEtcd(ctx, serviceName, cli, false, true, opts...)
}

func dialWithEtcd(ctx context.Context, serviceName string, cli *gooetcd.Client, insecure, withCtx bool, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	if cli == nil || cli.Client == nil {
		return nil, fmt.Errorf("etcd client is nil")
	}
	builder, err := resolver.NewBuilder(cli.Client)
	if err != nil {
		return nil, err
	}

	target := fmt.Sprintf("%s:///%s", builder.Scheme(), serviceName)
	dialOpts := []grpc.DialOption{
		grpc.WithResolvers(builder),
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                4 * time.Minute,
			Timeout:             15 * time.Second,
			PermitWithoutStream: true,
		}),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(MaxRecvMsgSize),
			grpc.MaxCallSendMsgSize(MaxSendMsgSize),
		),
		grpc.WithChainUnaryInterceptor(clientUnaryInterceptorLog()),
		grpc.WithChainStreamInterceptor(clientStreamInterceptorLog()),
	}
	if insecure {
		dialOpts = append([]grpc.DialOption{grpc.WithInsecure()}, dialOpts...)
	}
	dialOpts = append(dialOpts, opts...)

	if withCtx {
		return grpc.DialContext(ctx, target, dialOpts...)
	}
	return grpc.Dial(target, dialOpts...)
}
