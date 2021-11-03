package goo_grpc

import (
	goo_etcd "github.com/liqiongtao/googo.io/goo-etcd"
)

type Option struct {
	Name  string
	Value interface{}
}

var (
	__env         = "env"
	__projectName = "project-name"
	__serviceName = "service-name"
	__version     = "version"
	__etcd        = "etcd-client"
)

func EnvOption(env string) Option {
	return Option{Name: __env, Value: env}
}

func ProjectNameOption(projectName string) Option {
	return Option{Name: __projectName, Value: projectName}
}

func ServiceNameOption(serviceName string) Option {
	return Option{Name: __serviceName, Value: serviceName}
}

func VersionOption(version string) Option {
	return Option{Name: __version, Value: version}
}

func EtcdOption(cli *goo_etcd.Client) Option {
	return Option{Name: __etcd, Value: cli}
}
