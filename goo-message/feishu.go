package goo_message

import (
	"encoding/json"
	"errors"
	goo_http_request "github.com/liqiongtao/googo.io/goo-http-request"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
	"runtime"
	"sync"
)

var (
	__fieShuCH   chan struct{}
	__feiShuOnce sync.Once
)

func FeiShu(hookUrl string, text string) error {
	__feiShuOnce.Do(func() {
		__fieShuCH = make(chan struct{}, runtime.NumCPU()*2)
	})

	__fieShuCH <- struct{}{}
	defer func() { <-__fieShuCH }()

	data := map[string]interface{}{
		"msg_type": "text",
		"content": map[string]interface{}{
			"text": text,
		},
	}
	b, _ := json.Marshal(&data)

	buf, err := goo_http_request.PostJson(hookUrl, b)
	if err != nil {
		return err
	}

	rst, _ := goo_utils.Byte(buf).Params()
	if msg := rst.Get("StatusMessage").String(); msg != "success" {
		return errors.New(msg)
	}

	return nil
}
