package goo_log

type Level int

// 定义日志级别
const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
	PANIC
	FATAL
)

var (
	// 定义日志级别文本
	LevelText = map[Level]string{
		DEBUG: "DEBUG",
		INFO:  "INFO",
		WARN:  "WARN",
		ERROR: "ERROR",
		PANIC: "PANIC",
		FATAL: "FATAL",
	}
)
