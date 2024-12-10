package goo_message

import (
	"encoding/json"
	"errors"
	"fmt"
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
	b, err := json.Marshal(&data)
	if err != nil {
		fmt.Println("[goo-msg][1001]", text, err)
		return err
	}

	buf, err := goo_http_request.PostJson(hookUrl, b)
	if err != nil {
		fmt.Println("[goo-msg][1002]", text, err)
		return err
	}

	rst, err := goo_utils.Byte(buf).Params()
	if err != nil {
		fmt.Println("[goo-msg][1003]", text, err)
		return err
	}
	if msg := rst.Get("StatusMessage").String(); msg != "success" {
		fmt.Println("[goo-msg][1004]", text, msg)
		return errors.New(msg)
	}

	return nil
}
