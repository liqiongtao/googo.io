package goo_grpc

type Config struct {
	ENV         string
	ServiceName string
	Version     string
	Addr        string
	PProfEnable bool
	PProfAddr   string
}
