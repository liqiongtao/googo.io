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
	if cli == nil || cli.Client == nil {
		return nil, fmt.Errorf("etcd client is nil")
	}
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
				// 客户端在该时间内未收到任何数据时发送 ping
				// 建议设置比服务端的 Time (5分钟) 小一些，确保客户端的连接不会因为服务端的 idle 检测而断开
				Time: 4 * time.Minute,
				// ping 请求的超时时间
				// 建议设置比服务端的 Timeout (20秒) 小一些
				Timeout: 15 * time.Second,
				// 允许在没有活动流的情况下发送ping
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
	if cli == nil || cli.Client == nil {
		return nil, fmt.Errorf("etcd client is nil")
	}
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
				// 客户端在该时间内未收到任何数据时发送 ping
				// 建议设置比服务端的 Time (5分钟) 小一些，确保客户端的连接不会因为服务端的 idle 检测而断开
				Time: 4 * time.Minute,
				// ping 请求的超时时间
				// 建议设置比服务端的 Timeout (20秒) 小一些
				Timeout: 15 * time.Second,
				// 允许在没有活动流的情况下发送ping
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
