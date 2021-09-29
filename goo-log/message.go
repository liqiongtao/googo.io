package goo_log

import (
	"encoding/json"
	"time"
)

type Message struct {
	Level   Level
	Time    time.Time
	Content string
	Data    map[string]interface{}
}

func (msg *Message) WithField(field string, value interface{}) {
	msg.Data[field] = value
}

func (msg *Message) String() string {
	return string(msg.JSON())
}

func (msg *Message) JSON() []byte {
	data := map[string]interface{}{
		"level": LevelText[msg.Level],
		"time":  msg.Time.Format("2006-01-02 15:04:05"),
		"msg":   msg.Content,
	}
	for k, v := range msg.Data {
		data[k] = v
	}
	buf, _ := json.Marshal(&data)
	return buf
}
