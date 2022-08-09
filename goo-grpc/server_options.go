package goo_grpc

import (
	goo_etcd "github.com/liqiongtao/googo.io/goo-etcd"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"time"
)

var defaultServerOptions = serverOptions{
	AuthFunc:      nil,
	EtcdClient:    nil,
	Register2Etcd: false,
	ServerOptions: []grpc.ServerOption{
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     15 * time.Second,
			MaxConnectionAge:      30 * time.Second,
			MaxConnectionAgeGrace: 5 * time.Second,
			Time:                  5 * time.Second,
			Timeout:               1 * time.Second,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             5 * time.Second,
			PermitWithoutStream: true,
		}),
	},
}

// 定义配置项集合
type serverOptions struct {
	AuthFunc AuthFunc

	EtcdClient    *goo_etcd.Client
	Register2Etcd bool

	ServerOptions []grpc.ServerOption
}

// 定义配置项抽象
type ServerOption interface {
	apply(*serverOptions)
}

// 定义配置项方法实现
type funcOption struct {
	f func(*serverOptions)
}

func (f *funcOption) apply(opts *serverOptions) {
	f.f(opts)
}

func newFuncOption(f func(opts *serverOptions)) *funcOption {
	return &funcOption{f: f}
}

// 配置项 - 认证方法
func AuthFuncOption(authFunc AuthFunc) ServerOption {
	return newFuncOption(func(opts *serverOptions) {
		opts.AuthFunc = authFunc
	})
}

// 配置项 - grpc server options
func ServerOptions(opt ...grpc.ServerOption) ServerOption {
	return newFuncOption(func(opts *serverOptions) {
		opts.ServerOptions = append(opts.ServerOptions, opt...)
	})
}
