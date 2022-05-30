package goo_redis

import (
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
)

var __clients = map[string]*Redis{}

func Init(configs ...Config) {
	for _, conf := range configs {
		name := conf.Name
		if name == "" {
			name = "default"
		}
		__clients[name] = NewRedis(conf)
		if err := __clients[name].connect(); err != nil {
			goo_log.WithTag("goo-redis").Panic(err)
		}
		if conf.AutoPing {
			goo_utils.AsyncFunc(__clients[name].ping)
		}
	}
}

func Client(names ...string) *Redis {
	name := "default"
	if l := len(names); l > 0 {
		name = names[0]
	}
	if client, ok := __clients[name]; ok {
		return client
	}
	if l := len(__clients); l == 1 {
		for _, client := range __clients {
			return client
		}
	}
	goo_log.WithTag("goo-redis").Error("no default redis client")
	return nil
}
