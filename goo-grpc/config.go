package goo_grpc

type Config struct {
	// 服务名称
	ServiceName string `json:"service_name" yaml:"service_name"`

	// 服务地址
	Addr string `json:"addr" yaml:"addr"`
}
