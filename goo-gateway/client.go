package goo_gateway

import (
	"sync"

	goo_etcd "github.com/liqiongtao/googo.io/goo-etcd"
	goo_grpc "github.com/liqiongtao/googo.io/goo-grpc"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"google.golang.org/grpc"
)

var (
	__cc = map[string]*grpc.ClientConn{}
	__mu sync.Mutex
)

func Client(serviceName string) *grpc.ClientConn {
	__mu.Lock()
	defer __mu.Unlock()

	if cc, ok := __cc[serviceName]; ok && cc != nil {
		return cc
	}

	cc, err := goo_grpc.DialWithEtcd(serviceName, goo_etcd.Default())
	if err != nil || cc == nil {
		goo_log.WithTag("goo-gateway").WithField("serviceName", serviceName).Error(err)
		return nil
	}
	__cc[serviceName] = cc
	return cc
}
