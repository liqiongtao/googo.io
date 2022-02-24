package goo_log

var __log *Logger

func init() {
	__log = New(NewConsoleAdapter())
}

func Default() *Logger {
	return __log
}

func SetAdapter(adapter Adapter) {
	__log.SetAdapter(adapter)
}

func SetFileAdapterOptions(opts ...Option) {
	if adapter, ok := __log.adapter.(*FileAdapter); ok {
		adapter.SetOptions(opts...)
	}
}

func SetTrimPath(trimPaths ...string) {
	__log.SetTrimPath(trimPaths...)
}

func AddHook(fn func(msg Message)) {
	__log.AddHook(fn)
}

func WithTag(tags ...string) *Entry {
	return __log.WithTag(tags...)
}

func WithField(field string, value interface{}) *Entry {
	return __log.WithField(field, value)
}

func WithTrace() *Entry {
	return __log.WithTrace()
}

func Debug(v ...interface{}) {
	__log.Debug(v...)
}

func DebugF(format string, v ...interface{}) {
	__log.DebugF(format, v...)
}

func Info(v ...interface{}) {
	__log.Info(v...)
}

func InfoF(format string, v ...interface{}) {
	__log.InfoF(format, v...)
}

func Warn(v ...interface{}) {
	__log.Warn(v...)
}

func WarnF(format string, v ...interface{}) {
	__log.WarnF(format, v...)
}

func Error(v ...interface{}) {
	__log.Error(v...)
}

func ErrorF(format string, v ...interface{}) {
	__log.ErrorF(format, v...)
}

func Panic(v ...interface{}) {
	__log.Panic(v...)
}

func PanicF(format string, v ...interface{}) {
	__log.PanicF(format, v...)
}

func Fatal(v ...interface{}) {
	__log.Fatal(v...)
}

func FatalF(format string, v ...interface{}) {
	__log.FatalF(format, v...)
}
