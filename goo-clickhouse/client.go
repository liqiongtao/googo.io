package goo_clickhouse

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
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
	if conf.Driver == "" {
		conf.Driver = "clickhouse"
	}

	cli = &Client{Config: conf}

	if err = cli.connect(); err != nil {
		cli = nil
		return
	}

	if err = cli.DB.Ping(); err != nil {
		goo_log.WithTag("goo-clickhouse").Error(err)
		_ = cli.DB.Close()
		cli = nil
		return
	}

	return
}

func (cli *Client) connect() (err error) {
	addrs := []string{cli.Config.Addr}
	if cli.Config.AltHosts != "" {
		for _, h := range strings.Split(cli.Config.AltHosts, ",") {
			h = strings.TrimSpace(h)
			if h != "" {
				addrs = append(addrs, h)
			}
		}
	}

	opts := &clickhouse.Options{
		Addr: addrs,
		Auth: clickhouse.Auth{
			Database: cli.Config.Database,
			Username: cli.Config.User,
			Password: cli.Config.Password,
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		DialTimeout:      10 * time.Second,
		ReadTimeout:      time.Duration(cli.Config.ReadTimeout) * time.Second,
		ConnMaxLifetime:     time.Hour,
		MaxOpenConns:     20,
		MaxIdleConns:     5,
		ConnOpenStrategy: clickhouse.ConnOpenInOrder,
		Debug:            cli.Config.Debug,
	}

	cli.DB = clickhouse.OpenDB(opts)
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
		goo_log.WithTag("goo-clickhouse").
			WithField("err_code", exception.Code).
			WithField("stack_trace", exception.StackTrace).
			Error(exception.Message)
		return
	}

	goo_log.WithTag("goo-clickhouse").Error(err)
}
