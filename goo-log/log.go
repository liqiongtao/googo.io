package goo_log

import "sync"

var (
	__log  *Logger
	__once sync.Once
)

func Default() *Logger {
	__once.Do(func() {
		__log = NewConsoleLog()
	})
	return __log
}

func SetAdapter(adapter Adapter) {
	Default().SetAdapter(adapter)
}

func Sync() error {
	return Default().Sync()
}

func Close() error {
	return Default().Close()
}

func WithHook(fns ...func(msg *Message)) {
	Default().WithHook(fns...)
}

func WithTag(tags ...string) *Entry {
	return Default().WithTag(tags...)
}

func WithField(field string, value any) *Entry {
	return Default().WithField(field, value)
}

func WithTrace() *Entry {
	return Default().WithTrace()
}

func Debug(v ...any) {
	Default().Debug(v...)
}

func DebugF(format string, v ...any) {
	Default().DebugF(format, v...)
}

func Info(v ...any) {
	Default().Info(v...)
}

func InfoF(format string, v ...any) {
	Default().InfoF(format, v...)
}

func Warn(v ...any) {
	Default().Warn(v...)
}

func WarnF(format string, v ...any) {
	Default().WarnF(format, v...)
}

func Error(v ...any) {
	Default().Error(v...)
}

func ErrorF(format string, v ...any) {
	Default().ErrorF(format, v...)
}

func Panic(v ...any) {
	Default().Panic(v...)
}

func PanicF(format string, v ...any) {
	Default().PanicF(format, v...)
}

func Fatal(v ...any) {
	Default().Fatal(v...)
}

func FatalF(format string, v ...any) {
	Default().FatalF(format, v...)
}
