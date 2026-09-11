package goodb

import (
	"fmt"
	"strings"

	goolog "github.com/liqiongtao/googo.io/goo-log"
	"xorm.io/xorm/log"
)

type logger struct {
	showSQL  bool
	LogLevel log.LogLevel
	l        *goolog.Logger
}

func newLogger(logFilePath string) *logger {
	if logFilePath == "" {
		logFilePath = "logs/sql/"
	}
	return &logger{
		l: goolog.NewFileLog(goolog.FilePathOption(logFilePath)),
	}
}

func (l logger) Debug(v ...any) {
	if l.LogLevel > log.LOG_DEBUG {
		return
	}
	l.l.Debug(v...)
}

func (l logger) Debugf(format string, v ...any) {
	if l.LogLevel > log.LOG_DEBUG {
		return
	}
	l.l.Debug(fmt.Sprintf(format, v...))
}

func (l logger) Error(v ...any) {
	if l.LogLevel > log.LOG_ERR {
		return
	}
	l.l.Error(v...)
}

func (l logger) Errorf(format string, v ...any) {
	if l.LogLevel > log.LOG_ERR {
		return
	}
	l.l.Error(fmt.Sprintf(format, v...))
}

func (l logger) Info(v ...any) {
	if l.LogLevel > log.LOG_INFO {
		return
	}
	l.l.Info(v...)
}

func (l logger) Infof(format string, v ...any) {
	if l.LogLevel > log.LOG_INFO {
		return
	}
	if strings.Index(format, "PING DATABASE") != -1 {
		return
	}
	l.l.Info(fmt.Sprintf(format, v...))
}

func (l logger) Warn(v ...any) {
	if l.LogLevel > log.LOG_WARNING {
		return
	}
	l.l.Warn(v...)
}

func (l logger) Warnf(format string, v ...any) {
	if l.LogLevel > log.LOG_WARNING {
		return
	}
	l.l.Warn(fmt.Sprintf(format, v...))
}

func (l logger) Level() log.LogLevel {
	return l.LogLevel
}

func (l *logger) SetLevel(ll log.LogLevel) {
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

func (l *logger) Close() error {
	if l == nil || l.l == nil {
		return nil
	}
	return l.l.Close()
}
