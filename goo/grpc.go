package goo

import (
	"context"
	goo_grpc "github.com/liqiongtao/googo.io/goo-grpc"
)

func NewGRPCServer(cfg goo_grpc.Config) *goo_grpc.Server {
	return goo_grpc.New(cfg)
}

func GRPCContext(ctx *Context) context.Context {
	return goo_grpc.Context(ctx.Context)
}
