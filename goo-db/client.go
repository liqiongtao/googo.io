package goo_db

import (
	"context"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
)

var __clients = map[string]*Orm{}

func Init(ctx context.Context, configs ...Config) {
	for _, config := range configs {
		name := config.Name
		if name == "" {
			name = "default"
		}
		__clients[name] = New(ctx, config)
		if err := __clients[name].connect(); err != nil {
			continue
		}
		if config.AutoPing {
			goo_utils.AsyncFunc(__clients[name].ping)
		}
	}
}

func Client(names ...string) *Orm {
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

	goo_log.WithTag("goo-db").Error("no default db client")
	return nil
}
