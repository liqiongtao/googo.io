package goo_gateway

import (
	goo_etcd "github.com/liqiongtao/googo.io/goo-etcd"
	goo_grpc "github.com/liqiongtao/googo.io/goo-grpc"
	"google.golang.org/grpc"
	"sync"
)

var (
	__cc = map[string]*grpc.ClientConn{}
	__mu sync.Mutex
)

func Client(serviceName string) *grpc.ClientConn {
	__mu.Lock()
	defer __mu.Unlock()

	if _, ok := __cc[serviceName]; !ok {
		cc, _ := goo_grpc.DialWithEtcd(serviceName, goo_etcd.Default())
		__cc[serviceName] = cc
	}

	return __cc[serviceName]
}
