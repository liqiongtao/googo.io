package goo_grpc

type Config struct {
	// 服务名称
	ServiceName string `json:"service_name" yaml:"service_name"`

	// 对外开放地址
	ServiceEndpoint string `json:"service_endpoint" yaml:"service_endpoint"`

	// 监听地址
	Addr string `json:"addr" yaml:"addr"`
}
