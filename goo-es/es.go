package goo_es

import (
	"github.com/elastic/go-elasticsearch/v7"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"net"
	"net/http"
	"time"
)

var __client *client

func Init(conf Config) {
	__client, _ = New(conf)
}

func New(conf Config) (*client, error) {
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
		goo_log.WithTag("goo-es").WithField("config", cfg).Error(err)
		return nil, err
	}

	return &client{cli}, nil
}
