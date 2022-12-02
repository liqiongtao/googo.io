package goo_mongo

import (
	"context"
	"fmt"
	goo_cron "github.com/liqiongtao/googo.io/goo-cron"
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

	uri := fmt.Sprintf("mongodb://%s:%s@%s/%s", conf.User, conf.Password, conf.Addr, conf.Database)
	opts := options.Client().ApplyURI(uri)

	cli.Client, err = mongo.Connect(cli.ctx, opts)
	if err != nil {
		goo_log.WithTag("goo-mongo").Error(err)
		return
	}

	if err = cli.Ping(cli.ctx, readpref.Primary()); err != nil {
		goo_log.WithTag("goo-mongo").Error(err)
		return
	}

	if conf.AutoPing {
		goo_cron.SecondX(5, func() {
			if err := cli.Ping(cli.ctx, readpref.Primary()); err != nil {
				goo_log.WithTag("goo-mongo").Error(err)
			}
		})
	}

	return
}

func (cli *Client) WithContext(ctx context.Context) *Client {
	cli.ctx = ctx
	return cli
}

func (cli *Client) DB() *mongo.Database {
	return cli.Database(cli.conf.Database)
}
