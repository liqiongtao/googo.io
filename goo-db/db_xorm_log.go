package goo_db

import (
	"fmt"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"strings"
	"xorm.io/core"
)

type xormLogger struct {
	LogLevel core.LogLevel
	l        *goo_log.Logger
}

func newXormLogger(logFilePath, logFileName string) *xormLogger {
	return &xormLogger{
		l: goo_log.New(goo_log.NewFileAdapter(
			goo_log.FileDirnameOption(logFilePath),
			goo_log.FileTagOption(logFileName),
		)),
	}
}

func (l xormLogger) Debug(v ...interface{}) {
	l.l.Debug(v...)
}

func (l xormLogger) Debugf(format string, v ...interface{}) {
	l.l.Debug(fmt.Sprintf(format, v...))
}

func (l xormLogger) Error(v ...interface{}) {
	l.l.Error(v...)
}

func (l xormLogger) Errorf(format string, v ...interface{}) {
	l.l.Error(fmt.Sprintf(format, v...))
}

func (l xormLogger) Info(v ...interface{}) {
	l.l.Info(v...)
}

func (l xormLogger) Infof(format string, v ...interface{}) {
	if strings.Index(format, "PING DATABASE") != -1 {
		return
	}
	l.l.Info(fmt.Sprintf(format, v...))
}

func (l xormLogger) Warn(v ...interface{}) {
	l.l.Warn(v...)
}

func (l xormLogger) Warnf(format string, v ...interface{}) {
	l.l.Warn(fmt.Sprintf(format, v...))
}

func (l xormLogger) Level() core.LogLevel {
	return l.LogLevel
}

func (l xormLogger) SetLevel(ll core.LogLevel) {
	l.LogLevel = ll
}

func (l xormLogger) ShowSQL(show ...bool) {
}

func (l xormLogger) IsShowSQL() bool {
	return true
}
