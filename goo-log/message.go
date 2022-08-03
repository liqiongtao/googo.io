package goo_log

import (
	"encoding/json"
	"fmt"
	"time"
)

type Message struct {
	Level   Level
	Message []interface{}
	Trace   []string
	Time    time.Time
	Entry   *Entry
}

func (msg *Message) JSON() []byte {
	data := map[string]interface{}{}

	if l := len(msg.Entry.Data); l > 0 {
		for _, i := range msg.Entry.Data {
			data[i.Field] = i.Value
		}
	}

	{
		info := map[string]interface{}{}

		info["level"] = LevelText[msg.Level]
		info["datetime"] = msg.Time.Format("2006-01-02 15:04:05")

		if l := len(msg.Entry.Tags); l > 0 {
			info["tags"] = msg.Entry.Tags
		}

		if l := len(msg.Message); l > 0 {
			var arr []string
			for _, i := range msg.Message {
				arr = append(arr, fmt.Sprint(i))
			}
			info["message"] = arr
		}

		if l := len(msg.Trace); l > 0 {
			info["trace"] = msg.Trace
		}

		data["log"] = info
	}

	buf, _ := json.Marshal(&data)
	return buf
}
