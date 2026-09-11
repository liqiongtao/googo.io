package goomail

import (
	"bytes"
	"fmt"
	"strings"
)

type Message struct {
	Sender     string   // 发送者
	SenderName string   // 发送者名称
	Receivers  []string // 接受者
	Subject    string   // 主题
	Body       string   // 内容
}

func (m *Message) Html() []byte {
	var bf bytes.Buffer

	bf.WriteString("To: " + strings.Join(sanitizeHeaderValues(m.Receivers), ",\r\n "))
	bf.WriteString("\r\n")
	if m.SenderName != "" {
		bf.WriteString(fmt.Sprintf("From: \"%s\"<%s>", sanitizeHeaderValue(m.SenderName), sanitizeHeaderValue(m.Sender)))
	} else {
		bf.WriteString("From: " + sanitizeHeaderValue(m.Sender))
	}
	bf.WriteString("\r\n")
	bf.WriteString("Subject: " + sanitizeHeaderValue(m.Subject))
	bf.WriteString("\r\n")
	bf.WriteString("Content-Type:text/html;charset=utf-8")
	bf.WriteString("\r\n\r\n")
	bf.WriteString(m.Body)

	return bf.Bytes()
}

func sanitizeHeaderValue(v string) string {
	v = strings.ReplaceAll(v, "\r", "")
	v = strings.ReplaceAll(v, "\n", "")
	return v
}

func sanitizeHeaderValues(vs []string) []string {
	out := make([]string, len(vs))
	for i, v := range vs {
		out[i] = sanitizeHeaderValue(v)
	}
	return out
}
