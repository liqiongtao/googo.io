package goo_es

import (
	"net"
	"net/http"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
)

var __client *ESClient

func Init(conf Config) error {
	cli, err := New(conf)
	if err != nil {
		return err
	}
	__client = cli
	return nil
}

func Client() *ESClient {
	return __client
}

func New(conf Config) (*ESClient, error) {
	cfg := elasticsearch.Config{
		Addresses:         conf.Addresses,
		Username:          conf.User,
		Password:          conf.Password,
		EnableDebugLogger: conf.EnableLog,
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: 30 * time.Second,
			DialContext:           (&net.Dialer{Timeout: 30 * time.Second}).DialContext,
		},
	}

	cli, err := elasticsearch.NewClient(cfg)
	if err != nil {
		goo_log.WithTag("goo-es").
			WithField("addresses", conf.Addresses).
			WithField("user", conf.User).
			Error(err)
		return nil, err
	}

	return &ESClient{cli: cli}, nil
}
