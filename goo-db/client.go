package goo_db

import (
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
)

type Client struct {
	*xorm.EngineGroup
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
		return
	}

	if err = cli.Ping(); err != nil {
		goo_log.WithTag("goo-db").Error(err)
		return
	}

	cli.EngineGroup.ShowSQL(conf.LogModel)
	cli.EngineGroup.SetLogger(newLogger(conf.LogFilepath))
	cli.EngineGroup.SetMaxIdleConns(conf.MaxIdle)
	cli.EngineGroup.SetMaxOpenConns(conf.MaxOpen)
	if conf.MaxLifetime > 0 {
		cli.EngineGroup.SetConnMaxLifetime(time.Duration(conf.MaxLifetime) * time.Second)
	} else {
		cli.EngineGroup.SetConnMaxLifetime(600 * time.Second)
	}

	return
}
