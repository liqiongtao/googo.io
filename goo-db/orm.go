package goo_db

import (
	"context"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	goo_context "github.com/liqiongtao/googo.io/goo-context"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"time"
)

type Orm struct {
	*xorm.EngineGroup

	conf Config
	ctx  context.Context
}

func New(conf Config) *Orm {
	return &Orm{
		ctx:  goo_context.Cancel(),
		conf: conf,
	}
}

func (db *Orm) connect() (err error) {
	conns := []string{db.conf.Master}
	if n := len(db.conf.Slaves); n > 0 {
		conns = append(conns, db.conf.Slaves...)
	}

	db.EngineGroup, err = xorm.NewEngineGroup(db.conf.Driver, conns)
	if err != nil {
		goo_log.WithTag("goo-db").Error(err)
		return
	}

	var (
		logFilepath = "logs/sql/"
	)
	if db.conf.LogFilepath != "" {
		logFilepath = db.conf.LogFilepath
	}
	db.EngineGroup.SetLogger(newLogger(logFilepath))

	db.EngineGroup.ShowSQL(db.conf.LogModel)
	db.EngineGroup.SetMaxIdleConns(db.conf.MaxIdle)
	db.EngineGroup.SetMaxOpenConns(db.conf.MaxOpen)

	return
}

func (db *Orm) ping() {
	ti := time.NewTimer(time.Second)
	defer ti.Stop()

	for {
		select {
		case <-db.ctx.Done():
			return

		case <-ti.C:
			if err := db.EngineGroup.Ping(); err != nil {
				goo_log.WithTag("goo-db").Error(err)
			}
			ti.Reset(5 * time.Second)
		}
	}
}
