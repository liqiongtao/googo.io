package goo_log

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

func NewEntry(l *Logger) *Entry {
	return &Entry{
		l: l,
		msg: &Message{
			Data: map[string]interface{}{},
		},
	}
}

type Entry struct {
	l   *Logger
	msg *Message
}

func (entry *Entry) WithField(field string, value interface{}) *Entry {
	entry.msg.WithField(field, value)
	return entry
}

func (entry *Entry) Debug(v ...interface{}) {
	entry.output(DEBUG, v...)
}

func (entry *Entry) Info(v ...interface{}) {
	entry.output(INFO, v...)
}

func (entry *Entry) Warn(v ...interface{}) {
	entry.output(WARN, v...)
}

func (entry *Entry) Error(v ...interface{}) {
	entry.output(ERROR, v...)
}

func (entry *Entry) Panic(v ...interface{}) {
	entry.output(PANIC, v...)
}

func (entry *Entry) Fatal(v ...interface{}) {
	entry.output(FATAL, v...)
}

func (entry *Entry) output(level Level, v ...interface{}) {
	entry.msg.Level = level
	entry.msg.Time = time.Now()
	entry.msg.Content = fmt.Sprint(v...)

	if level >= ERROR {
		entry.msg.WithField("trace", entry.trace())
	}

	for _, hook := range entry.l.hooks {
		go hook(*entry.msg)
	}

	if entry.l.adapter != nil {
		entry.l.adapter.Write(entry.msg)
	}
}

func (entry *Entry) trace() (arr []string) {
	arr = []string{}
	ll := len(entry.l.trimPaths)
	for i := 3; i < 16; i++ {
		_, file, line, _ := runtime.Caller(i)
		if file == "" ||
			strings.Index(file, "runtime/") > 0 ||
			strings.Index(file, "src/") > 0 ||
			strings.Index(file, "pkg/mod/") > 0 ||
			strings.Index(file, "vendor/") > 0 {
			continue
		}
		if ll > 0 {
			for _, trimPath := range entry.l.trimPaths {
				file = strings.Replace(file, trimPath, "", -1)
			}
		}
		arr = append([]string{fmt.Sprintf("%s %dL", file, line)}, arr...)
	}
	return
}
