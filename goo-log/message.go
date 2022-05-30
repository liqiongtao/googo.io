package goo_log

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
	"time"
)

type Message struct {
	Level   Level
	Message []interface{}
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

	data["log_level"] = LevelText[msg.Level]
	data["log_datetime"] = msg.Time.Format("2006-01-02 15:04:05")

	if l := len(msg.Entry.Tags); l > 0 {
		data["log_tags"] = msg.Entry.Tags
	}

	if l := len(msg.Message); l > 0 {
		var arr []string
		for _, i := range msg.Message {
			arr = append(arr, fmt.Sprint(i))
		}
		data["log_content"] = arr
	}

	if msg.Level >= WARN {
		data["log_trace"] = msg.trace()
	}

	buf, _ := json.Marshal(&data)
	return buf
}

func (msg *Message) trace() (arr []string) {
	arr = []string{}

	for i := 2; i < 16; i++ {
		_, file, line, _ := runtime.Caller(i)
		if file == "" ||
			strings.Contains(file, ".pb.go") ||
			strings.Contains(file, "runtime/") ||
			strings.Contains(file, "src/testing") ||
			strings.Contains(file, "pkg/mod/") ||
			strings.Contains(file, "vendor/") {
			continue
		}
		arr = append(arr, fmt.Sprintf("%s %dL", msg.prettyFile(file), line))
	}

	return
}

func (msg *Message) prettyFile(file string) string {
	index := strings.LastIndex(file, "/")
	if index < 0 {
		return file
	}

	index2 := strings.LastIndex(file[:index], "/")
	if index2 < 0 {
		return file[index+1:]
	}

	return file[index2+1:]
}
