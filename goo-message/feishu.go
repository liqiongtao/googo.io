package goo_message

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"sync"

	goo_request "github.com/liqiongtao/googo.io/goo-request"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
)

var (
	__fieShuCH   chan struct{}
	__feiShuOnce sync.Once
	__feiShuMu   sync.Mutex
)

func FeiShu(hookUrl string, text string) error {
	if text == "" {
		return nil
	}

	__feiShuOnce.Do(func() {
		__fieShuCH = make(chan struct{}, runtime.NumCPU())
	})

	// 控制并发
	__fieShuCH <- struct{}{}
	defer func() { <-__fieShuCH }()

	data := map[string]any{
		"msg_type": "text",
		"content": map[string]any{
			"text": text,
		},
	}
	b, err := json.Marshal(&data)
	if err != nil {
		fmt.Println("[goo-msg][1001]", text, err)
		return err
	}

	buf, err := goo_request.PostJson(hookUrl, b)
	if err != nil {
		fmt.Println("[goo-msg][1002]", text, err)
		return err
	}
	if len(buf) == 0 {
		return errors.New("empty feishu response")
	}
	if bytes.Contains(buf, []byte("服务异常，请联系")) {
		return errors.New("feishu service error")
	}

	rst, err := goo_utils.Byte(buf).Params()
	if err != nil {
		fmt.Println("[goo-msg][1003]", text, string(buf), err)
		return err
	}

	msg := rst.Get("msg").String()
	switch msg {
	case "success":
		return nil
	case "too many request":
		return errors.New("feishu rate limited")
	default:
		fmt.Println("[goo-msg][1004]", text, string(buf))
		return errors.New(msg)
	}
}
