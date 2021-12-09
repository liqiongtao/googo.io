package goo_log

func New(adapter Adapter) *Logger {
	return &Logger{
		adapter:   adapter,
		hooks:     []func(msg Message){},
		trimPaths: []string{},
	}
}

type Logger struct {
	adapter   Adapter
	hooks     []func(msg Message)
	trimPaths []string
}

func (l *Logger) SetAdapter(adapter Adapter) {
	l.adapter = adapter
}

func (l *Logger) SetFileAdapterOptions(opts ...Option) {
	if adapter, ok := l.adapter.(*FileAdapter); ok {
		adapter.SetOptions(opts...)
	}
}

func (l *Logger) AddHook(fn func(msg Message)) {
	l.hooks = append(l.hooks, fn)
}

func (l *Logger) SetTrimPath(trimPaths ...string) {
	l.trimPaths = append(l.trimPaths, trimPaths...)
}

func (l *Logger) WithTag(tags ...string) *Entry {
	return NewEntry(l).WithTag(tags...)
}

func (l *Logger) WithField(field string, value interface{}) *Entry {
	return NewEntry(l).WithField(field, value)
}

func (l *Logger) Debug(v ...interface{}) {
	NewEntry(l).Debug(v...)
}

func (l *Logger) Info(v ...interface{}) {
	NewEntry(l).Info(v...)
}

func (l *Logger) Warn(v ...interface{}) {
	NewEntry(l).Warn(v...)
}

func (l *Logger) Error(v ...interface{}) {
	NewEntry(l).Error(v...)
}

func (l *Logger) Panic(v ...interface{}) {
	NewEntry(l).Panic(v...)
}

func (l *Logger) Fatal(v ...interface{}) {
	NewEntry(l).Fatal(v...)
}
