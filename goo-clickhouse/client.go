package goo_clickhouse

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/ClickHouse/clickhouse-go"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
)

type Client struct {
	Config
	*sql.DB
}

func New(conf Config) (cli *Client, err error) {
	if conf.ReadTimeout == 0 {
		conf.ReadTimeout = 10
	}
	if conf.WriteTimeout == 0 {
		conf.WriteTimeout = 20
	}

	cli = &Client{Config: conf}

	if err = cli.connect(); err != nil {
		return
	}

	return
}

func (cli *Client) connect() (err error) {
	dns := fmt.Sprintf("tcp://%s?username=%s&password=%s&database=%s&read_timeout=%d&write_timeout=%d&alt_hosts=%s&debug=%v",
		cli.Config.Addr, cli.Config.User, cli.Config.Password, cli.Config.Database,
		cli.Config.ReadTimeout, cli.Config.WriteTimeout, cli.Config.AltHosts, cli.Config.Debug)
	cli.DB, err = sql.Open(cli.Config.Driver, dns)
	if err != nil {
		goo_log.WithTag("goo-clickhouse").Error(err)
	}
	return
}

func (cli *Client) ping() {
	if cli.DB == nil {
		return
	}

	err := cli.DB.Ping()
	if err == nil {
		return
	}

	var exception *clickhouse.Exception
	if errors.As(err, &exception) {
		goo_log.WithTag("goo-clickhouse").WithField("err_code", exception.Code).WithField("stack_trace", exception.StackTrace).Error(exception.Message)
		return
	}

	goo_log.WithTag("goo-clickhouse").Error(err)
}
