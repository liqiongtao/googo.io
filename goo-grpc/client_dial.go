package goo_grpc

import (
	"context"
	"fmt"
	goo_etcd "github.com/liqiongtao/googo.io/goo-etcd"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
)

func Dial(ctx context.Context, addr string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	opts = append(opts,
		grpc.WithInsecure(),
	)
	return grpc.DialContext(ctx, addr, opts...)
}

func DialWithEtcd(ctx context.Context, serviceName string, cli *goo_etcd.Client) (*grpc.ClientConn, error) {
	b := NewResolverEtcd(ctx, cli)
	serviceName = fmt.Sprintf("%s:///%s", b.Scheme(), serviceName)

	opts := []grpc.DialOption{
		grpc.WithResolvers(b),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)),
		grpc.WithInsecure(),
	}

	return grpc.DialContext(ctx, serviceName, opts...)
}
