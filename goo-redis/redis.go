package goo_redis

import (
	goo_log "github.com/liqiongtao/googo.io/goo-log"
)

var __clients = map[string]*Client{}

func Init(configs ...Config) (err error) {
	for _, conf := range configs {
		name := conf.Name
		if name == "" {
			name = "default"
		}

		cli, e := New(conf)
		if e != nil {
			return e
		}
		__clients[name] = cli
	}

	return
}

func GetClient(names ...string) *Client {
	name := "default"
	if l := len(names); l > 0 {
		name = names[0]
	}

	if cli, ok := __clients[name]; ok {
		return cli
	}

	if l := len(__clients); l == 1 {
		for _, cli := range __clients {
			return cli
		}
	}

	goo_log.WithTag("goo-redis").Error("no default redis client")

	return nil
}

func Default() *Client {
	if cli, ok := __clients["default"]; ok {
		return cli
	}

	if l := len(__clients); l == 1 {
		for _, cli := range __clients {
			return cli
		}
	}

	goo_log.WithTag("goo-redis").Error("no default redis client")

	return nil
}
