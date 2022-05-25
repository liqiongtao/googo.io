package goo_db

import (
	"context"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"time"
)

func NewXOrmAdapter(ctx context.Context, config Config) *XOrm {
	return &XOrm{
		ctx:    ctx,
		config: config,
	}
}

type XOrm struct {
	ctx    context.Context
	config Config
	*xorm.EngineGroup
}

func (db *XOrm) connect() (err error) {
	conns := []string{db.config.Master}
	if n := len(db.config.Slaves); n > 0 {
		conns = append(conns, db.config.Slaves...)
	}

	db.EngineGroup, err = xorm.NewEngineGroup(db.config.Driver, conns)
	if err != nil {
		return
	}

	var (
		logFilepath = "logs/sql/"
	)
	if db.config.LogFilepath != "" {
		logFilepath = db.config.LogFilepath
	}
	db.EngineGroup.SetLogger(newXormLogger(logFilepath))

	db.EngineGroup.ShowSQL(db.config.LogModel)
	db.EngineGroup.SetMaxIdleConns(db.config.MaxIdle)
	db.EngineGroup.SetMaxOpenConns(db.config.MaxOpen)

	return
}

func (db *XOrm) ping() {
	dur := 3 * time.Second
	ti := time.NewTimer(dur)
	defer ti.Stop()
	for {
		select {
		case <-db.ctx.Done():
			return
		case <-ti.C:
			if err := db.EngineGroup.Ping(); err != nil {
				goo_log.WithTag("goo-db").Error(err)
			}
			ti.Reset(dur)
		}
	}
}
