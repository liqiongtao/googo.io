package goo_grpc

import (
	"context"
	"fmt"
	goo_etcd "github.com/liqiongtao/googo.io/goo-etcd"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/keepalive"
	"time"
)

func DialWithEtcd(serviceName string, cli *goo_etcd.Client) (*grpc.ClientConn, error) {
	builder, err := resolver.NewBuilder(cli.Client)
	if err != nil {
		return nil, err
	}

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithResolvers(builder),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                10 * time.Second,
			Timeout:             100 * time.Millisecond,
			PermitWithoutStream: true,
		}),
		grpc.WithChainUnaryInterceptor(clientUnaryInterceptorLog()),
		grpc.WithChainStreamInterceptor(clientStreamInterceptorLog()),
	}

	return grpc.Dial(builder.Scheme()+":///"+serviceName, opts...)
}

func DialContextWithEtcd(ctx context.Context, serviceName string, cli *goo_etcd.Client) (*grpc.ClientConn, error) {
	builder, err := resolver.NewBuilder(cli.Client)
	if err != nil {
		return nil, err
	}

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithResolvers(builder),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                10 * time.Second,
			Timeout:             100 * time.Millisecond,
			PermitWithoutStream: true,
		}),
		grpc.WithChainUnaryInterceptor(clientUnaryInterceptorLog()),
		grpc.WithChainStreamInterceptor(clientStreamInterceptorLog()),
	}

	return grpc.DialContext(ctx, builder.Scheme()+":///"+serviceName, opts...)
}
