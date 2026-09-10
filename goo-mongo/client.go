package goo_mongo

import (
	"context"
	"fmt"
	"net/url"
	"time"

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
	timeout := conf.Timeout
	if timeout <= 0 {
		timeout = 10
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	cli = &Client{conf: conf, ctx: context.Background()}

	var uri string
	if conf.User != "" {
		authSource := conf.AuthSource
		if authSource == "" {
			authSource = "admin"
		}
		userInfo := url.UserPassword(conf.User, conf.Password)
		uri = fmt.Sprintf("mongodb://%s@%s/%s?authSource=%s",
			userInfo.String(), conf.Addr, url.PathEscape(conf.Database), url.QueryEscape(authSource))
	} else {
		uri = fmt.Sprintf("mongodb://%s/%s",
			conf.Addr, url.PathEscape(conf.Database))
	}
	opts := options.Client().ApplyURI(uri)

	cli.Client, err = mongo.Connect(ctx, opts)
	if err != nil {
		goo_log.WithTag("goo-mongo").Error(err)
		cli = nil
		return
	}

	if err = cli.Ping(ctx, readpref.Primary()); err != nil {
		goo_log.WithTag("goo-mongo").Error(err)
		_ = cli.Disconnect(context.Background())
		cli = nil
		return
	}

	return
}

func (cli *Client) WithContext(ctx context.Context) *Client {
	if ctx == nil {
		ctx = context.Background()
	}
	return &Client{Client: cli.Client, conf: cli.conf, ctx: ctx}
}

// Context 返回 WithContext 绑定的上下文，供 Find/Insert 等操作使用。
func (cli *Client) Context() context.Context {
	if cli == nil || cli.ctx == nil {
		return context.Background()
	}
	return cli.ctx
}

func (cli *Client) DB() *mongo.Database {
	return cli.Database(cli.conf.Database)
}
