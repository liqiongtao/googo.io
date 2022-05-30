package goo_log

import (
	"fmt"
	"log"
	"os"
	"time"
)

type Entry struct {
	Tags []string
	Data []DataField
	msg  *Message
	l    *Logger
}

type DataField struct {
	Field string
	Value interface{}
}

func NewEntry(l *Logger) *Entry {
	return &Entry{l: l}
}

func (entry *Entry) WithTag(tags ...string) *Entry {
	entry.Tags = append(entry.Tags, tags...)
	return entry
}

func (entry *Entry) WithField(field string, value interface{}) *Entry {
	entry.Data = append(entry.Data, DataField{Field: field, Value: value})
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
}

func (entry *Entry) PanicF(format string, v ...interface{}) {
	entry.output(PANIC, fmt.Sprintf(format, v...))
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
	entry.msg = &Message{
		Level:   level,
		Message: v,
		Time:    time.Now(),
		Entry:   entry,
	}

	for _, fn := range entry.l.hooks {
		go entry.hookHandler(fn)
	}

	if entry.l.adapter != nil {
		entry.l.adapter.Write(entry.msg)
	}
}

func (entry *Entry) hookHandler(fn func(msg *Message)) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()

	fn(entry.msg)
}
