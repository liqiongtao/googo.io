package goo_etcd

import (
	"context"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

var __cli *Client

func CLI() *Client {
	return __cli
}

func Init(endpoints []string) {
	var err error
	__cli, err = New(endpoints)
	if err != nil {
		goo_log.Panic(err.Error())
	}
}

func New(endpoints []string) (cli *Client, err error) {
	cli = &Client{ctx: context.Background()}
	cli.Client, err = clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		goo_log.Error(err.Error())
	}
	return
}

func Set(key, val string, ttl int64) (err error) {
	return __cli.Set(key, val, ttl)
}

func SetWithKeepAlive(key, val string, ttl int64) (err error) {
	return __cli.SetWithKeepAlive(key, val, ttl)
}

func Get(key string) string {
	return __cli.Get(key)
}

func GetMap(key string) (data map[string]string) {
	return __cli.GetMap(key)
}

func GetArray(key string) (data []string) {
	return __cli.GetArray(key)
}

func Watch(key string, fn func(arr []string)) {
	__cli.Watch(key, fn)
}

func Del(key string) error {
	return __cli.Del(key)
}
