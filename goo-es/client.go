package goo_es

import (
	"context"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
)

type client struct {
	cli *elasticsearch.Client
}

func (c *client) Client() *elasticsearch.Client {
	return c.cli
}

func (c *client) log() *goo_log.Entry {
	return goo_log.WithTag("goo-es")
}

func (c *client) exec(req esapi.Request) (*esapi.Response, error) {
	res, err := req.Do(context.TODO(), c.cli)
	if err != nil {
		c.log().Error(err)
		return nil, err
	}
	return res, nil
}
