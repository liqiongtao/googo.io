package goo_clickhouse

import (
	"database/sql"
	"fmt"
	"github.com/ClickHouse/clickhouse-go"
	"github.com/liqiongtao/googo.io/goo"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
)

type client struct {
	conf Config
	db   *sql.DB
}

func New(conf Config) (cli *client, err error) {
	if conf.ReadTimeout == 0 {
		conf.ReadTimeout = 10
	}
	if conf.WriteTimeout == 0 {
		conf.WriteTimeout = 20
	}

	cli = &client{conf: conf}

	if err = cli.connect(); err != nil {
		return
	}

	if !conf.AutoPing || conf.PingDuration == 0 {
		return
	}

	goo.Crond(conf.PingDuration, __client.ping)
	return
}

func (cli *client) connect() (err error) {
	dns := fmt.Sprintf("tcp://%s?username=%s&password=%s&database=%s&read_timeout=%d&write_timeout=%d&alt_hosts=%s&debug=%v",
		cli.conf.Addr, cli.conf.User, cli.conf.Password, cli.conf.Database,
		cli.conf.ReadTimeout, cli.conf.WriteTimeout, cli.conf.AltHosts, cli.conf.Debug)
	cli.db, err = sql.Open(cli.conf.Driver, dns)
	if err != nil {
		goo_log.WithTag("goo-clickhouse").Error(err)
	}
	return
}

func (cli *client) ping() {
	if cli.db == nil {
		return
	}

	err := cli.db.Ping()
	if err == nil {
		return
	}

	if exception, ok := err.(*clickhouse.Exception); ok {
		goo_log.WithTag("goo-clickhouse").WithField("err_code", exception.Code).WithField("stack_trace", exception.StackTrace).Error(exception.Message)
		return
	}

	goo_log.WithTag("goo-clickhouse").Error(err)
}
