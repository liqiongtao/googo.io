package goo_grpc

import "time"

type Config struct {
	// 服务名称
	ServiceName string `json:"service_name" yaml:"service_name"`

	// 对外开放地址
	ServiceEndpoint string `json:"service_endpoint" yaml:"service_endpoint"`

	// 监听地址
	Addr string `json:"addr" yaml:"addr"`

	// 超时时间(单位秒)
	KeepaliveTime    time.Duration `json:"keepalive_time" yaml:"keepalive_time"`
	KeepaliveTimeout time.Duration `json:"keepalive_timeout" yaml:"keepalive_timeout"`

	// ping包最短等待时间(单位秒)
	EnforcementPolicyMinTime time.Duration `json:"enforcement_policy_min_time" yaml:"enforcement_policy_min_time"`
}
