package goo

import (
	"encoding/json"
)

type Response struct {
	Code    int32         `json:"code"`
	Message string        `json:"message"`
	Data    interface{}   `json:"data"`
	Errors  []interface{} `json:"-"`
}

func (rsp *Response) String() string {
	buf, err := json.Marshal(rsp)
	if err != nil {
		return err.Error()
	}
	return string(buf)
}

func Success(data interface{}) *Response {
	if data == nil {
		data = map[string]interface{}{}
	}
	return &Response{
		Code:    0,
		Message: "ok",
		Data:    data,
	}
}

func Error(code int32, message string, v ...interface{}) *Response {
	return &Response{
		Code:    code,
		Message: message,
		Data:    map[string]string{},
		Errors:  v,
	}
}
