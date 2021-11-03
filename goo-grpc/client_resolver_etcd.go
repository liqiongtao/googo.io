package goo_grpc

import (
	"context"
	goo_etcd "github.com/liqiongtao/googo.io/goo-etcd"
	"google.golang.org/grpc/resolver"
)

type ResolverEtcd struct {
	ctx context.Context
	cli *goo_etcd.Client
}

func NewResolverEtcd(ctx context.Context, cli *goo_etcd.Client) *ResolverEtcd {
	return &ResolverEtcd{ctx, cli}
}

func (b *ResolverEtcd) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	b.cli.Watch(target.Endpoint, func(arr []string) {
		var addresses []resolver.Address
		for _, addr := range arr {
			addresses = append(addresses, resolver.Address{Addr: addr})
		}
		cc.UpdateState(resolver.State{Addresses: addresses})
	})

	r := &Resolver{}
	r.ResolveNow(resolver.ResolveNowOptions{})

	return r, nil
}

func (b *ResolverEtcd) Scheme() string {
	return "etcd"
}
