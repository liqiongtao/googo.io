package goo_etcd

import (
	"context"
	"crypto/tls"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"go.etcd.io/etcd/client/pkg/v3/transport"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

var __cli *Client

func CLI() *Client {
	return __cli
}

func Init(cfg Config) {
	var err error
	__cli, err = New(cfg)
	if err != nil {
		goo_log.Panic(err.Error())
	}
}

func New(cfg Config) (cli *Client, err error) {
	cli = &Client{ctx: context.Background()}

	config := clientv3.Config{
		Endpoints:   cfg.Endpoints,
		DialTimeout: 5 * time.Second,
		// 该日志配置，用于屏蔽原生日志输出
		LogConfig: &zap.Config{
			Level:    zap.NewAtomicLevelAt(zapcore.ErrorLevel),
			Encoding: "json",
		},
	}

	if cfg.User != "" {
		config.Username = cfg.User
	}
	if cfg.Password != "" {
		config.Password = cfg.Password
	}

	if cfg.TLS != nil {
		tlsInfo := &transport.TLSInfo{
			CertFile:      cfg.TLS.CertFile,
			KeyFile:       cfg.TLS.KeyFile,
			TrustedCAFile: cfg.TLS.CAFile,
		}
		var clientConfig *tls.Config
		clientConfig, err = tlsInfo.ClientConfig()
		if err != nil {
			goo_log.Error(err.Error())
			return
		}
		config.TLS = clientConfig
	}

	cli.Client, err = clientv3.New(config)
	if err != nil {
		goo_log.Error(err.Error())
	}
	return
}

func Set(key, val string, ttl int64) (err error) {
	return __cli.Set(key, val, ttl)
}

func SetWithKeepAlive(key, val string, ttl int64) (err error) {
	return __cli.SetWithKeepAlive(key, val, ttl)
}

func Get(key string) string {
	return __cli.Get(key)
}

func GetMap(key string) (data map[string]string) {
	return __cli.GetMap(key)
}

func GetArray(key string) (data []string) {
	return __cli.GetArray(key)
}

func Watch(key string, fn func(arr []string)) {
	__cli.Watch(key, fn)
}

func Del(key string) error {
	return __cli.Del(key)
}
