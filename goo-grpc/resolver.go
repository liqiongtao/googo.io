package goo_grpc

import (
	"context"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
	"google.golang.org/grpc/resolver"
)

type Resolver struct {
	ctx    context.Context
	target resolver.Target
	cc     resolver.ClientConn
}

// 监测并更新grpc节点变化
func (r *Resolver) Watch(ch chan []resolver.Address) {
	goo_utils.AsyncFunc(func() {
		for {
			select {
			case <-r.ctx.Done():
				return
			case addresses := <-ch:
				r.cc.UpdateState(resolver.State{Addresses: addresses})
				goo_log.WithField("endpoint", r.target.Endpoint).
					WithField("addresses", addresses).
					Debug("[grpc-nodes:update]")
			}
		}
	})
}

func (r *Resolver) ResolveNow(o resolver.ResolveNowOptions) {
}

func (r *Resolver) Close() {
}
