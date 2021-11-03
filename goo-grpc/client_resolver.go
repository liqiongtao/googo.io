package goo_grpc

import (
	"google.golang.org/grpc/resolver"
)

type Resolver struct {
}

func (r *Resolver) ResolveNow(o resolver.ResolveNowOptions) {
}

func (r *Resolver) Close() {
}
