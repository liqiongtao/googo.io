package goo_grpc

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/resolver"
)

// 返回客户端对象
func Dial(serviceName string, ch chan []resolver.Address, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	b := NewResolverBuilder(ch)
	opts = append(opts,
		grpc.WithResolvers(b),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)),
		grpc.WithInsecure(),
	)
	return grpc.Dial(fmt.Sprintf("%s:///%s", b.Scheme(), serviceName), opts...)
}
