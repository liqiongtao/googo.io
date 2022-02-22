package goo_log

import (
	"encoding/json"
	"time"
)

type Message struct {
	Level   Level
	Tags    []string
	Time    time.Time
	Content string
	Data    map[string]interface{}
}

func (msg *Message) WithTag(tags ...string) {
	msg.Tags = append(msg.Tags, tags...)
}

func (msg *Message) WithField(field string, value interface{}) {
	msg.Data[field] = value
}

func (msg *Message) String() string {
	return string(msg.JSON())
}

func (msg *Message) JSON() []byte {
	data := map[string]interface{}{
		"__level": LevelText[msg.Level],
		"__time":  msg.Time.Format("2006-01-02 15:04:05"),
	}

	if l := len(msg.Tags); l > 0 {
		data["__tag"] = msg.Tags
	}

	if msg.Content != "" {
		data["__msg"] = msg.Content
	}

	for k, v := range msg.Data {
		data[k] = v
	}

	buf, _ := json.Marshal(&data)
	return buf
}
