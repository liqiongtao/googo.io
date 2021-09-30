package goo_db

import (
	"context"
	goo_log "googo.io/goo-log"
	goo_utils "googo.io/goo-utils"
)

var __clients = map[string]DB{}

func Init(ctx context.Context, configs ...Config) {
	for _, config := range configs {
		name := config.Name
		if name == "" {
			name = "default"
		}
		__clients[name] = NewXOrmAdapter(ctx, config)
		if err := __clients[name].connect(); err != nil {
			panic(err.Error())
		}
		if config.AutoPing {
			goo_utils.AsyncFunc(__clients[name].ping)
		}
	}
}

func Client(names ...string) DB {
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
	goo_log.Error("no default db client")
	return nil
}

func XOrmClient(names ...string) *XOrm {
	client := Client(names...)
	if client == nil {
		return nil
	}
	return client.(*XOrm)
}
