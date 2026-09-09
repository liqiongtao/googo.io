package goo_log

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

type Entry struct {
	Tags  []string
	Data  []DataField
	Trace []string
	msg   *Message
	l     *Logger
}

type DataField struct {
	Field string
	Value interface{}
}

func NewEntry(l *Logger) *Entry {
	return &Entry{l: l}
}

// clone 拷贝 Tags/Data/Trace，保证 With* 不会污染原 Entry。
func (entry *Entry) clone() *Entry {
	e := &Entry{l: entry.l}
	if n := len(entry.Tags); n > 0 {
		e.Tags = make([]string, n)
		copy(e.Tags, entry.Tags)
	}
	if n := len(entry.Data); n > 0 {
		e.Data = make([]DataField, n)
		copy(e.Data, entry.Data)
	}
	if n := len(entry.Trace); n > 0 {
		e.Trace = make([]string, n)
		copy(e.Trace, entry.Trace)
	}
	return e
}

func (entry *Entry) WithTag(tags ...string) *Entry {
	if len(tags) == 0 {
		return entry
	}
	e := entry.clone()
	e.Tags = append(e.Tags, tags...)
	return e
}

func (entry *Entry) WithField(field string, value interface{}) *Entry {
	e := entry.clone()
	e.Data = append(e.Data, DataField{Field: field, Value: value})
	return e
}

func (entry *Entry) WithTrace() *Entry {
	e := entry.clone()
	e.Trace = e.trace()
	return e
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
	_ = entry.l.Sync()
	os.Exit(1)
}

func (entry *Entry) FatalF(format string, v ...interface{}) {
	entry.output(FATAL, fmt.Sprintf(format, v...))
	_ = entry.l.Sync()
	os.Exit(1)
}

func (entry *Entry) output(level Level, v ...interface{}) {
	// 拷贝后再写，避免同一 Entry 并发打日志时互相覆盖 msg/Trace
	e := entry.clone()
	e.msg = &Message{
		Level:   level,
		Message: v,
		Time:    time.Now(),
		Entry:   e,
	}

	// 保留 WithTrace() 预置的栈；WARN+ 且未预置时自动采集
	if level >= WARN && len(e.Trace) == 0 {
		e.Trace = e.trace()
	}

	_, hooks := e.l.snapshot()
	for _, fn := range hooks {
		e.hookHandler(fn)
	}

	adapter, _ := e.l.snapshot()
	if adapter != nil {
		adapter.Write(e.msg)
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

// runtime.Caller 仅能获取非 goroutine 的信息
func (entry *Entry) trace() (arr []string) {
	arr = []string{}

	for i := 3; i < 16; i++ {
		_, file, line, _ := runtime.Caller(i)
		if file == "" {
			continue
		}
		if strings.Contains(file, ".pb.go") ||
			strings.Contains(file, "runtime/") ||
			(!strings.Contains(file, "googo.io") &&
				(strings.Contains(file, "src/") || strings.Contains(file, "pkg/mod/") || strings.Contains(file, "vendor/"))) ||
			strings.Contains(file, "goo-log") {
			continue
		}
		arr = append(arr, fmt.Sprintf("%s %dL", entry.prettyFile(file), line))
	}

	return
}

func (entry *Entry) prettyFile(file string) string {
	var (
		index  int
		index2 int
	)

	if index = strings.LastIndex(file, "src/test/"); index >= 0 {
		return file[index+9:]
	}
	if index = strings.LastIndex(file, "src/"); index >= 0 {
		return file[index+4:]
	}
	if index = strings.LastIndex(file, "pkg/mod/"); index >= 0 {
		return file[index+8:]
	}
	if index = strings.LastIndex(file, "vendor/"); index >= 0 {
		return file[index+7:]
	}

	if index = strings.LastIndex(file, "/"); index < 0 {
		return file
	}
	if index2 = strings.LastIndex(file[:index], "/"); index2 < 0 {
		return file[index+1:]
	}
	return file[index2+1:]
}
