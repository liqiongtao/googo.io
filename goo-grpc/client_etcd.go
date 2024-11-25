package goo_grpc

import (
	"context"
	"fmt"
	"time"

	goo_etcd "github.com/liqiongtao/googo.io/goo-etcd"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

func DialWithEtcd(serviceName string, cli *goo_etcd.Client, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	builder, err := resolver.NewBuilder(cli.Client)
	if err != nil {
		return nil, err
	}

	// 使用 URL 格式构建 target
	target := fmt.Sprintf("%s:///%s", builder.Scheme(), serviceName)

	opts = append(
		[]grpc.DialOption{
			grpc.WithInsecure(),
			grpc.WithResolvers(builder),
			grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`),
			grpc.WithKeepaliveParams(keepalive.ClientParameters{
				Time:                30 * time.Second,
				Timeout:             10 * time.Second,
				PermitWithoutStream: true,
			}),
			grpc.WithDefaultCallOptions(
				grpc.MaxCallRecvMsgSize(MaxRecvMsgSize),
				grpc.MaxCallSendMsgSize(MaxSendMsgSize),
			),
			grpc.WithChainUnaryInterceptor(clientUnaryInterceptorLog()),
			grpc.WithChainStreamInterceptor(clientStreamInterceptorLog()),
		},
		opts...,
	)

	return grpc.Dial(target, opts...)
}

func DialContextWithEtcd(ctx context.Context, serviceName string, cli *goo_etcd.Client, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	builder, err := resolver.NewBuilder(cli.Client)
	if err != nil {
		return nil, err
	}

	// 使用 URL 格式构建 target
	target := fmt.Sprintf("%s:///%s", builder.Scheme(), serviceName)

	opts = append(
		[]grpc.DialOption{
			grpc.WithInsecure(),
			grpc.WithResolvers(builder),
			grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`),
			grpc.WithKeepaliveParams(keepalive.ClientParameters{
				Time:                30 * time.Second,
				Timeout:             10 * time.Second,
				PermitWithoutStream: true,
			}),
			grpc.WithDefaultCallOptions(
				grpc.MaxCallRecvMsgSize(MaxRecvMsgSize),
				grpc.MaxCallSendMsgSize(MaxSendMsgSize),
			),
			grpc.WithChainUnaryInterceptor(clientUnaryInterceptorLog()),
			grpc.WithChainStreamInterceptor(clientStreamInterceptorLog()),
		},
		opts...,
	)

	return grpc.DialContext(ctx, target, opts...)
}
