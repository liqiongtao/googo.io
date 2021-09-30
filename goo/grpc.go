package goo

import (
	"context"
	goo_grpc "googo.io/goo-grpc"
)

func NewGRPCServer() *goo_grpc.Server {
	return goo_grpc.New()
}

func GRPCContext(ctx *Context) context.Context {
	return goo_grpc.Context(ctx.Context)
}
