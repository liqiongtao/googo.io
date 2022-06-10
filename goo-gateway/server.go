package goo_gateway

import "github.com/liqiongtao/googo.io/goo"

func Server(conf Config, opts ...goo.Option) {
	opts = append([]goo.Option{
		goo.ServerNameOption(conf.ServerName),
		goo.EnvOption(conf.Env),
	}, opts...)

	s := goo.NewServer(opts...)

	s.POST("/:service/:method", goo.Handler(gateway{conf: conf}))

	s.Run(conf.Addr)
}
