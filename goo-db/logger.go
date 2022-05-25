package goo_db

import (
	"fmt"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"strings"
	"xorm.io/core"
)

type logger struct {
	LogLevel core.LogLevel
	l        *goo_log.Logger
}

func newLogger(logFilePath string) *logger {
	return &logger{
		l: goo_log.NewFileLog(goo_log.FilePathOption(logFilePath)),
	}
}

func (l logger) Debug(v ...interface{}) {
	l.l.Debug(v...)
}

func (l logger) Debugf(format string, v ...interface{}) {
	l.l.Debug(fmt.Sprintf(format, v...))
}

func (l logger) Error(v ...interface{}) {
	l.l.Error(v...)
}

func (l logger) Errorf(format string, v ...interface{}) {
	l.l.Error(fmt.Sprintf(format, v...))
}

func (l logger) Info(v ...interface{}) {
	l.l.Info(v...)
}

func (l logger) Infof(format string, v ...interface{}) {
	if strings.Index(format, "PING DATABASE") != -1 {
		return
	}
	l.l.Info(fmt.Sprintf(format, v...))
}

func (l logger) Warn(v ...interface{}) {
	l.l.Warn(v...)
}

func (l logger) Warnf(format string, v ...interface{}) {
	l.l.Warn(fmt.Sprintf(format, v...))
}

func (l logger) Level() core.LogLevel {
	return l.LogLevel
}

func (l logger) SetLevel(ll core.LogLevel) {
	l.LogLevel = ll
}

func (l logger) ShowSQL(show ...bool) {
}

func (l logger) IsShowSQL() bool {
	return true
}
