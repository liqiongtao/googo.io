package main

import (
	"context"
	"fmt"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"time"
)

func Register2Etcd(conf Config) error {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{""},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		goo_log.WithTag("goo-grpc").Error(err)
		return err
	}

	m, err := endpoints.NewManager(cli, conf.ServiceName)
	if err != nil {
		goo_log.WithTag("goo-grpc").Error(err)
		return err
	}

	key := fmt.Sprintf("%s/%s", conf.ServiceName, conf.Addr)
	ep := endpoints.Endpoint{Addr: conf.Addr}

	if err := m.AddEndpoint(context.TODO(), key, ep); err != nil {
		goo_log.WithTag("goo-grpc").Error(err)
		return err
	}

	return nil
}
