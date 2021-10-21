package goo_grpc

import (
	"context"
	"google.golang.org/grpc/resolver"
)

type ResolverBuilder struct {
	ctx context.Context
	ch  chan []resolver.Address
}

func NewResolverBuilder(ctx context.Context, ch chan []resolver.Address) *ResolverBuilder {
	return &ResolverBuilder{ctx, ch}
}

func (b *ResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &Resolver{
		ctx:    b.ctx,
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
