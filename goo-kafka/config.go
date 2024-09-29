package goo_kafka

import goo_redis "github.com/liqiongtao/googo.io/goo-redis"

type Config struct {
	User              string   `json:"user" yaml:"user"`
	Password          string   `json:"password" yaml:"password"`
	Addrs             []string `json:"addrs" yaml:"addrs"`
	Timeout           int      `json:"timeout" yaml:"timeout"`
	HeartbeatInterval int      `json:"heartbeat_interval" yaml:"heartbeat_interval"`
	SessionTimeout    int      `json:"session_timeout" yaml:"session_timeout"`
	RebalanceTimeout  int      `json:"rebalance_timeout" yaml:"rebalance_timeout"`
	RedisConfig       goo_redis.Config
}
