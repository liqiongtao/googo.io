package goo_redis

import (
	"sync"

	goo_log "github.com/liqiongtao/googo.io/goo-log"
)

var (
	__clients = map[string]*Client{}
	__mu      sync.RWMutex
)

func Init(configs ...Config) (err error) {
	added := make([]string, 0, len(configs))
	for _, conf := range configs {
		name := conf.Name
		if name == "" {
			name = "default"
		}

		cli, e := New(conf)
		if e != nil {
			__mu.Lock()
			for _, n := range added {
				if c := __clients[n]; c != nil {
					_ = c.Close()
				}
				delete(__clients, n)
			}
			__mu.Unlock()
			return e
		}
		__mu.Lock()
		__clients[name] = cli
		__mu.Unlock()
		added = append(added, name)
	}

	return
}

func GetClient(names ...string) *Client {
	name := "default"
	if l := len(names); l > 0 {
		name = names[0]
	}

	__mu.RLock()
	defer __mu.RUnlock()

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
	__mu.RLock()
	defer __mu.RUnlock()

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
