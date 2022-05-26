package goo_etcd

type Config struct {
	Endpoints []string `json:"endpoints" yaml:"endpoints"`
	Username  string   `json:"username" yaml:"username"`
	Password  string   `json:"password" yaml:"password"`
	TLS       *TLS     `json:"tls" yaml:"tls"`
}

type TLS struct {
	CertFile string `json:"cert_file" yaml:"cert_file"`
	KeyFile  string `json:"key_file" yaml:"key_file"`
	CAFile   string `json:"ca_file" yaml:"ca_file"`
}
