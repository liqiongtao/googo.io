package goo_log

import "sync"

type Logger struct {
	mu      sync.RWMutex
	hooks   []func(msg *Message)
	adapter Adapter
}

func New(adapter Adapter) *Logger {
	return &Logger{
		adapter: adapter,
	}
}

func (l *Logger) SetAdapter(adapter Adapter) {
	l.mu.Lock()
	old := l.adapter
	l.adapter = adapter
	l.mu.Unlock()

	// 替换带后台 worker 的适配器时，避免泄漏
	if old != nil && old != adapter {
		if c, ok := old.(interface{ Close() error }); ok {
			_ = c.Close()
		}
	}
}

func (l *Logger) WithHook(fns ...func(msg *Message)) {
	l.mu.Lock()
	l.hooks = append(l.hooks, fns...)
	l.mu.Unlock()
}

func (l *Logger) snapshot() (Adapter, []func(msg *Message)) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	hooks := make([]func(msg *Message), len(l.hooks))
	copy(hooks, l.hooks)
	return l.adapter, hooks
}

func (l *Logger) Sync() error {
	adapter, _ := l.snapshot()
	if s, ok := adapter.(interface{ Sync() error }); ok {
		return s.Sync()
	}
	return nil
}

func (l *Logger) Close() error {
	adapter, _ := l.snapshot()
	if c, ok := adapter.(interface{ Close() error }); ok {
		return c.Close()
	}
	return nil
}

func (l *Logger) WithTag(tags ...string) *Entry {
	return NewEntry(l).WithTag(tags...)
}

func (l *Logger) WithField(field string, value interface{}) *Entry {
	return NewEntry(l).WithField(field, value)
}

func (l *Logger) WithTrace() *Entry {
	return NewEntry(l).WithTrace()
}

func (l *Logger) Debug(v ...interface{}) {
	NewEntry(l).Debug(v...)
}

func (l *Logger) DebugF(format string, v ...interface{}) {
	NewEntry(l).DebugF(format, v...)
}

func (l *Logger) Info(v ...interface{}) {
	NewEntry(l).Info(v...)
}

func (l *Logger) InfoF(format string, v ...interface{}) {
	NewEntry(l).InfoF(format, v...)
}

func (l *Logger) Warn(v ...interface{}) {
	NewEntry(l).Warn(v...)
}

func (l *Logger) WarnF(format string, v ...interface{}) {
	NewEntry(l).WarnF(format, v...)
}

func (l *Logger) Error(v ...interface{}) {
	NewEntry(l).Error(v...)
}

func (l *Logger) ErrorF(format string, v ...interface{}) {
	NewEntry(l).ErrorF(format, v...)
}

func (l *Logger) Panic(v ...interface{}) {
	NewEntry(l).Panic(v...)
}

func (l *Logger) PanicF(format string, v ...interface{}) {
	NewEntry(l).PanicF(format, v...)
}

func (l *Logger) Fatal(v ...interface{}) {
	NewEntry(l).Fatal(v...)
}

func (l *Logger) FatalF(format string, v ...interface{}) {
	NewEntry(l).FatalF(format, v...)
}
