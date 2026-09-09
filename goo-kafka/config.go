package goo_kafka

import goo_redis "github.com/liqiongtao/googo.io/goo-redis"

type Config struct {
	User              string           `json:"user" yaml:"user"`
	Password          string           `json:"password" yaml:"password"`
	Addrs             []string         `json:"addrs" yaml:"addrs"`
	Timeout           int              `json:"timeout" yaml:"timeout"`
	HeartbeatInterval int              `json:"heartbeat_interval" yaml:"heartbeat_interval"`
	SessionTimeout    int              `json:"session_timeout" yaml:"session_timeout"`
	RebalanceTimeout  int              `json:"rebalance_timeout" yaml:"rebalance_timeout"`
	ClientID          string           `json:"client_id" yaml:"client_id"`     // 可选；未配置则用随机 UUID
	InstanceId        string           `json:"instance_id" yaml:"instance_id"` // 静态成员 ID；未配置则不启用
	RedisConfig       goo_redis.Config `json:"redis_config" yaml:"redis_config"`
}
