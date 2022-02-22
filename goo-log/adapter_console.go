package goo_log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

// 控制台日志适配器
type ConsoleAdapter struct {
}

func NewConsoleAdapter() *ConsoleAdapter {
	return new(ConsoleAdapter)
}

func (ca ConsoleAdapter) Write(msg *Message) {
	var buf bytes.Buffer
	buf.WriteString(msg.Time.Format("2006-01-02 15:04:05"))
	buf.WriteString(" ")
	buf.WriteString(colors[msg.Level](fmt.Sprintf("%-5s", LevelText[msg.Level])))
	buf.WriteString(" ")
	if l := len(msg.Tags); l > 0 {
		buf.WriteString("[" + strings.Join(msg.Tags, "][") + "]")
		buf.WriteString(" ")
	}
	if msg.Content != "" {
		buf.WriteString(msg.Content)
		buf.WriteString(" ")
	}
	if l := len(msg.Data); l > 0 {
		if b, err := json.Marshal(&msg.Data); err == nil {
			buf.Write(b)
		}
	}
	ca.writer().Write(append(buf.Bytes(), '\n'))
}

func (ca ConsoleAdapter) writer() io.Writer {
	return os.Stdout
}

type brush func(string) string

// 定义日志级别颜色
var colors = map[Level]brush{
	DEBUG: newBrush("1;34"), // blue
	INFO:  newBrush("1;32"), // green
	WARN:  newBrush("1;33"), // yellow
	ERROR: newBrush("1;31"), // red
	PANIC: newBrush("1;37"), // white
	FATAL: newBrush("1;35"), // magenta
}

func newBrush(color string) brush {
	pre := "\033["
	reset := "\033[0m"
	return func(text string) string {
		return pre + color + "m" + text + reset
	}
}
