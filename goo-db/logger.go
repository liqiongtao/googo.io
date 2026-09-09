package goo_db

import (
	"fmt"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"strings"
	"xorm.io/core"
)

type logger struct {
	showSQL  bool
	LogLevel core.LogLevel
	l        *goo_log.Logger
}

func newLogger(logFilePath string) *logger {
	if logFilePath == "" {
		logFilePath = "logs/sql/"
	}
	return &logger{
		l: goo_log.NewFileLog(goo_log.FilePathOption(logFilePath)),
	}
}

func (l logger) Debug(v ...interface{}) {
	if l.LogLevel > core.LOG_DEBUG {
		return
	}
	l.l.Debug(v...)
}

func (l logger) Debugf(format string, v ...interface{}) {
	if l.LogLevel > core.LOG_DEBUG {
		return
	}
	l.l.Debug(fmt.Sprintf(format, v...))
}

func (l logger) Error(v ...interface{}) {
	if l.LogLevel > core.LOG_ERR {
		return
	}
	l.l.Error(v...)
}

func (l logger) Errorf(format string, v ...interface{}) {
	if l.LogLevel > core.LOG_ERR {
		return
	}
	l.l.Error(fmt.Sprintf(format, v...))
}

func (l logger) Info(v ...interface{}) {
	if l.LogLevel > core.LOG_INFO {
		return
	}
	l.l.Info(v...)
}

func (l logger) Infof(format string, v ...interface{}) {
	if l.LogLevel > core.LOG_INFO {
		return
	}
	if strings.Index(format, "PING DATABASE") != -1 {
		return
	}
	l.l.Info(fmt.Sprintf(format, v...))
}

func (l logger) Warn(v ...interface{}) {
	if l.LogLevel > core.LOG_WARNING {
		return
	}
	l.l.Warn(v...)
}

func (l logger) Warnf(format string, v ...interface{}) {
	if l.LogLevel > core.LOG_WARNING {
		return
	}
	l.l.Warn(fmt.Sprintf(format, v...))
}

func (l logger) Level() core.LogLevel {
	return l.LogLevel
}

func (l *logger) SetLevel(ll core.LogLevel) {
	l.LogLevel = ll
}

func (l *logger) ShowSQL(show ...bool) {
	if len(show) == 0 {
		l.showSQL = true
		return
	}
	l.showSQL = show[0]
}

func (l logger) IsShowSQL() bool {
	return l.showSQL
}
