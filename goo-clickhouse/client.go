package goo_clickhouse

import (
	"database/sql"
	"fmt"
	"github.com/ClickHouse/clickhouse-go"
	"github.com/liqiongtao/googo.io/goo"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
)

type client struct {
	config Config
	db     *sql.DB
}

func New(config Config) (cli *client, err error) {
	if config.ReadTimeout == 0 {
		config.ReadTimeout = 10
	}
	if config.WriteTimeout == 0 {
		config.WriteTimeout = 20
	}

	cli = &client{config: config}

	if err = cli.connect(); err != nil {
		return
	}

	if !config.AutoPing || config.PingDuration == 0 {
		return
	}

	goo.Crond(config.PingDuration, __client.ping)
	return
}

func (cli *client) connect() (err error) {
	dns := fmt.Sprintf("tcp://%s?username=%s&password=%s&database=%s&read_timeout=%d&write_timeout=%d&alt_hosts=%s&debug=%v",
		cli.config.Addr, cli.config.User, cli.config.Password, cli.config.Database,
		cli.config.ReadTimeout, cli.config.WriteTimeout, cli.config.AltHosts, cli.config.Debug)
	cli.db, err = sql.Open(cli.config.Driver, dns)
	if err != nil {
		goo_log.Error(err.Error())
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
		goo_log.WithField("err_code", exception.Code).WithField("stack_trace", exception.StackTrace).Error(exception.Message)
		return
	}

	goo_log.Error(err.Error())
}
