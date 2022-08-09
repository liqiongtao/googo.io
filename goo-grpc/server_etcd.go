package goo_grpc

import goo_etcd "github.com/liqiongtao/googo.io/goo-etcd"

func (s *Server) Register2Etcd(cli *goo_etcd.Client) *Server {
	defaultServerOptions.Register2Etcd = true
	defaultServerOptions.EtcdClient = cli
	return s
}
