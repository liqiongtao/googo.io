package goo_grpc

import "time"

type Config struct {
	// 服务名称
	ServiceName string `json:"service_name" yaml:"service_name"`

	// 对外开放地址
	ServiceEndpoint string `json:"service_endpoint" yaml:"service_endpoint"`

	// 监听地址
	Addr string `json:"addr" yaml:"addr"`

	// 最大空闲链接时间，即空闲链接在配置的时间内，未收到新的心跳和请求，则会将链接关闭，比向客户端发送一个GoAway；
	// 空闲链接的定义：最近未完成的RPC数变为0 的时间，或链接建立以来；
	// 默认是无穷
	MaxConnectionIdle time.Duration `json:"max_connection_idle" yaml:"max_connection_idle"`
	// 最长链接时间，当stream超过这个时间会发一个GoAway；为了防止短时间内发送大量的GoAway 会根据 MaxConnectionAge 时间间隔随机+/- 10%
	// 默认是无穷
	MaxConnectionAge time.Duration `json:"max_connection_age" yaml:"max_connection_age"`
	// 是对MaxConnectionAge 的一个补充，超过了最长链接时间后延长的时间
	// 默认是无穷
	MaxConnectionAgeGrace time.Duration `json:"max_connection_age_grace" yaml:"max_connection_age_grace"`

	// 服务端在设定的时间范围内未收到客户端任何活动，例如stream在时间内未收到数据信息，则会发送ping 信息检查链接是否可用；
	// 及时发现及时重试；
	// 当设置值小于1秒时，会被强制设置成1秒
	KeepaliveTime time.Duration `json:"keepalive_time" yaml:"keepalive_time"`
	// 服务端发送ping请求后，等待配置的时间，若客户端在这个时间内未有任何响应则将该链接关闭回收
	// 默认是20秒
	KeepaliveTimeout time.Duration `json:"keepalive_timeout" yaml:"keepalive_timeout"`

	// MinTime 是客户端在发送 keepalive ping 之前应等待的最短时间；
	// 即两个keepalive ping 之间的最小间隔，若小于这个间隔，则会关闭与客户端的链接
	// 默认是5分钟
	EnforcementPolicyMinTime time.Duration `json:"enforcement_policy_min_time" yaml:"enforcement_policy_min_time"`
}
