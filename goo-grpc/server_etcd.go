package goo_grpc

import goo_etcd "github.com/liqiongtao/googo.io/goo-etcd"

func (s *Server) Register2Etcd(cli *goo_etcd.Client) *Server {
	s.opts.Register2Etcd = true
	s.opts.EtcdClient = cli
	return s
}
