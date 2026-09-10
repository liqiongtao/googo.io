package goo_db

import (
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"xorm.io/xorm"
)

type Client struct {
	*xorm.EngineGroup
	sqlLogger *logger
}

func New(conf Config) (cli *Client, err error) {
	conns := []string{conf.Master}
	if n := len(conf.Slaves); n > 0 {
		conns = append(conns, conf.Slaves...)
	}

	cli = &Client{}

	cli.EngineGroup, err = xorm.NewEngineGroup(conf.Driver, conns)
	if err != nil {
		goo_log.WithTag("goo-db").Error(err)
		cli = nil
		return
	}

	if err = cli.Ping(); err != nil {
		goo_log.WithTag("goo-db").Error(err)
		_ = cli.Close()
		cli = nil
		return
	}

	if conf.LogModel {
		cli.sqlLogger = newLogger(conf.LogFilepath)
		cli.EngineGroup.SetLogger(cli.sqlLogger)
		cli.EngineGroup.ShowSQL(true)
	}

	maxIdle := conf.MaxIdle
	if maxIdle <= 0 {
		maxIdle = 10
	}
	maxOpen := conf.MaxOpen
	if maxOpen <= 0 {
		maxOpen = 100
	}
	cli.EngineGroup.SetMaxIdleConns(maxIdle)
	cli.EngineGroup.SetMaxOpenConns(maxOpen)
	if conf.MaxLifetime > 0 {
		cli.EngineGroup.SetConnMaxLifetime(time.Duration(conf.MaxLifetime) * time.Second)
	} else {
		cli.EngineGroup.SetConnMaxLifetime(600 * time.Second)
	}

	return
}

func (cli *Client) Close() error {
	if cli == nil {
		return nil
	}
	if cli.sqlLogger != nil {
		_ = cli.sqlLogger.Close()
		cli.sqlLogger = nil
	}
	if cli.EngineGroup != nil {
		return cli.EngineGroup.Close()
	}
	return nil
}
