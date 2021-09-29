package goo_redis

import (
	"context"
	goo_log "googo.io/goo-log"
	goo_utils "googo.io/goo-utils"
)

var __clients = map[string]*Redis{}

func Init(ctx context.Context, configs ...Config) {
	for _, config := range configs {
		name := config.Name
		if name == "" {
			name = "default"
		}
		__clients[name] = NewRedis(ctx, config)
		if err := __clients[name].connect(); err != nil {
			panic(err.Error())
		}
		if config.AutoPing {
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
	goo_log.Error("no default redis client")
	return nil
}
