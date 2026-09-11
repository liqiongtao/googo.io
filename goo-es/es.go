package gooes

import (
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	goolog "github.com/liqiongtao/googo.io/goo-log"
)

var (
	__client *ESClient
	__mu     sync.RWMutex
)

func Init(conf Config) error {
	cli, err := New(conf)
	if err != nil {
		return err
	}
	__mu.Lock()
	__client = cli
	__mu.Unlock()
	return nil
}

func Client() *ESClient {
	__mu.RLock()
	defer __mu.RUnlock()
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
		goolog.WithTag("goo-es").
			WithField("addresses", conf.Addresses).
			WithField("user", conf.User).
			Error(err)
		return nil, err
	}

	return &ESClient{cli: cli}, nil
}
