package goo_grpc

import (
	goo_etcd "github.com/liqiongtao/googo.io/goo-etcd"
)

func (s *Server) Register2ETCD(cli *goo_etcd.Client) error {
	return cli.RegisterService(s.conf.ServiceName, s.conf.Addr)
}
