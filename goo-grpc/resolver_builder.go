package goo_grpc

import (
	"google.golang.org/grpc/resolver"
)

type ResolverBuilder struct {
	ch chan []resolver.Address
}

func NewResolverBuilder(ch chan []resolver.Address) *ResolverBuilder {
	return &ResolverBuilder{ch}
}

func (b *ResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &Resolver{
		target: target,
		cc:     cc,
	}

	r.Watch(b.ch)

	r.ResolveNow(resolver.ResolveNowOptions{})

	return r, nil
}

func (b *ResolverBuilder) Scheme() string {
	return "goo"
}
