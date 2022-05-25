package goo_db

import (
	"context"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"time"
)

type Orm struct {
	*xorm.EngineGroup

	config Config
	ctx    context.Context
}

func New(ctx context.Context, config Config) *Orm {
	return &Orm{
		ctx:    ctx,
		config: config,
	}
}

func (db *Orm) connect() (err error) {
	conns := []string{db.config.Master}
	if n := len(db.config.Slaves); n > 0 {
		conns = append(conns, db.config.Slaves...)
	}

	db.EngineGroup, err = xorm.NewEngineGroup(db.config.Driver, conns)
	if err != nil {
		goo_log.WithTag("goo-db").Error(err)
		return
	}

	var (
		logFilepath = "logs/sql/"
	)
	if db.config.LogFilepath != "" {
		logFilepath = db.config.LogFilepath
	}
	db.EngineGroup.SetLogger(newLogger(logFilepath))

	db.EngineGroup.ShowSQL(db.config.LogModel)
	db.EngineGroup.SetMaxIdleConns(db.config.MaxIdle)
	db.EngineGroup.SetMaxOpenConns(db.config.MaxOpen)

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
