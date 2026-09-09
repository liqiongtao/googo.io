package goo_mongo

import (
	"context"
	"fmt"
	"net/url"

	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Client struct {
	*mongo.Client
	conf Config
	ctx  context.Context
}

func New(conf Config) (cli *Client, err error) {
	cli = &Client{conf: conf, ctx: context.TODO()}

	userInfo := url.UserPassword(conf.User, conf.Password)
	uri := fmt.Sprintf("mongodb://%s@%s/%s?authSource=%s",
		userInfo.String(), conf.Addr, url.PathEscape(conf.Database), url.QueryEscape(conf.Database))
	opts := options.Client().ApplyURI(uri)

	cli.Client, err = mongo.Connect(cli.ctx, opts)
	if err != nil {
		goo_log.WithTag("goo-mongo").Error(err)
		cli = nil
		return
	}

	if err = cli.Ping(cli.ctx, readpref.Primary()); err != nil {
		goo_log.WithTag("goo-mongo").Error(err)
		_ = cli.Disconnect(cli.ctx)
		cli = nil
		return
	}

	return
}

func (cli *Client) WithContext(ctx context.Context) *Client {
	return &Client{Client: cli.Client, conf: cli.conf, ctx: ctx}
}

func (cli *Client) DB() *mongo.Database {
	return cli.Database(cli.conf.Database)
}
