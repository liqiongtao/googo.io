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
			// 最大空闲链接时间，即空闲链接在配置的时间内，未收到新的心跳和请求，则会将链接关闭，比向客户端发送一个GoAway；
			// 空闲链接的定义：最近未完成的RPC数变为0 的时间，或链接建立以来；
			// 默认是无穷
			MaxConnectionIdle: 5 * time.Minute,
			// 最长链接时间，当stream超过这个时间会发一个GoAway；为了防止短时间内发送大量的GoAway 会根据 MaxConnectionAge 时间间隔随机+/- 10%
			// 默认是无穷
			MaxConnectionAge: 60 * time.Second,
			// 是对MaxConnectionAge 的一个补充，超过了最长链接时间后延长的时间
			// 默认是无穷
			MaxConnectionAgeGrace: 60 * time.Second,
			// 服务端在设定的时间范围内未收到客户端任何活动，例如stream在时间内未收到数据信息，则会发送ping 信息检查链接是否可用；
			// 及时发现及时重试；
			// 当设置值小于1秒时，会被强制设置成1秒
			Time: 5 * time.Second,
			// 服务端发送ping请求后，等待配置的时间，若客户端在这个时间内未有任何响应则将该链接关闭回收
			// 默认是20秒
			Timeout: 20 * time.Second,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			// MinTime 是客户端在发送 keepalive ping 之前应等待的最短时间；
			// 即两个keepalive ping 之间的最小间隔，若小于这个间隔，则会关闭与客户端的链接
			// 默认是5分钟
			MinTime: 60 * time.Second,
			// 没有 active stream, 也允许 ping
			// 如果为 false，并且客户端在没有活动流时发送 ping，服务器将发送 GoAway 并关闭连接
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
