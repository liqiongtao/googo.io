package goo_db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
	goo_cron "github.com/liqiongtao/googo.io/goo-cron"
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

	if conf.AutoPing {
		goo_cron.Default().SecondX(5, func() {
			if err := cli.Ping(); err != nil {
				goo_log.WithTag("goo-db").Error(err)
			}
		}).Start()
	}

	return
}
