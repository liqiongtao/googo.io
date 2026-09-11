package googrpc

import gooetcd "github.com/liqiongtao/googo.io/goo-etcd"

func (s *Server) Register2Etcd(cli *gooetcd.Client) *Server {
	s.opts.Register2Etcd = true
	s.opts.EtcdClient = cli
	return s
}
