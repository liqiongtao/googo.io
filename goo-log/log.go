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

func WithField(field string, value interface{}) *Entry {
	return __log.WithField(field, value)
}

func Debug(v ...interface{}) {
	__log.Debug(v...)
}

func Info(v ...interface{}) {
	__log.Info(v...)
}

func Warn(v ...interface{}) {
	__log.Warn(v...)
}

func Error(v ...interface{}) {
	__log.Error(v...)
}

func Panic(v ...interface{}) {
	__log.Panic(v...)
}

func Fatal(v ...interface{}) {
	__log.Fatal(v...)
}
