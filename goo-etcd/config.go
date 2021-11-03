package goo_etcd

type Config struct {
	User      string
	Password  string
	TLS       *TLS
	Endpoints []string
}

type TLS struct {
	CertFile string
	KeyFile  string
	CAFile   string
}
