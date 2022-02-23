package goo_log

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

func NewEntry(l *Logger) *Entry {
	return &Entry{
		l: l,
		msg: &Message{
			Tags: []string{},
			Data: map[string]interface{}{},
		},
	}
}

type Entry struct {
	l   *Logger
	msg *Message
}

func (entry *Entry) WithTag(tags ...string) *Entry {
	entry.msg.WithTag(tags...)
	return entry
}

func (entry *Entry) WithField(field string, value interface{}) *Entry {
	entry.msg.WithField(field, value)
	return entry
}

func (entry *Entry) Debug(v ...interface{}) {
	entry.output(DEBUG, v...)
}

func (entry *Entry) DebugF(format string, v ...interface{}) {
	entry.output(DEBUG, fmt.Sprintf(format, v...))
}

func (entry *Entry) Info(v ...interface{}) {
	entry.output(INFO, v...)
}

func (entry *Entry) InfoF(format string, v ...interface{}) {
	entry.output(INFO, fmt.Sprintf(format, v...))
}

func (entry *Entry) Warn(v ...interface{}) {
	entry.output(WARN, v...)
}

func (entry *Entry) WarnF(format string, v ...interface{}) {
	entry.output(WARN, fmt.Sprintf(format, v...))
}

func (entry *Entry) Error(v ...interface{}) {
	entry.output(ERROR, v...)
}

func (entry *Entry) ErrorF(format string, v ...interface{}) {
	entry.output(ERROR, fmt.Sprintf(format, v...))
}

func (entry *Entry) Panic(v ...interface{}) {
	entry.output(PANIC, v...)
	panic(fmt.Sprint(v...))
}

func (entry *Entry) PanicF(format string, v ...interface{}) {
	entry.output(PANIC, fmt.Sprintf(format, v...))
	panic(fmt.Sprintf(format, v...))
}

func (entry *Entry) Fatal(v ...interface{}) {
	entry.output(FATAL, v...)
	os.Exit(1)
}

func (entry *Entry) FatalF(format string, v ...interface{}) {
	entry.output(FATAL, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func (entry *Entry) output(level Level, v ...interface{}) {
	entry.msg.Level = level
	entry.msg.Time = time.Now()

	var arr []string
	for _, i := range v {
		arr = append(arr, fmt.Sprint(i))
	}
	entry.msg.Content = strings.Join(arr, " ")

	for _, trimPath := range entry.l.trimPaths {
		entry.msg.Content = strings.Replace(entry.msg.Content, trimPath, "", -1)
	}

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
		if file == "" {
			continue
		}
		if strings.Contains(file, "runtime/") ||
			strings.Contains(file, "src/") ||
			strings.Contains(file, ".pb.go") ||
			strings.Contains(file, "vendor/") ||
			(strings.Contains(file, "pkg/mod/") && !strings.Contains(file, "googo.io")) {
			continue
		}
		if ll > 0 {
			for _, trimPath := range entry.l.trimPaths {
				file = strings.Replace(file, trimPath, "", -1)
			}
		}
		arr = append(arr, fmt.Sprintf("%s %dL", file, line))
	}
	return
}
