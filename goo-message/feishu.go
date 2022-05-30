package goo_message

import (
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

func FieShu(hookUrl string, text string) error {
	__feiShuOnce.Do(func() {
		__fieShuCH = make(chan struct{}, runtime.NumCPU()*2)
	})

	__fieShuCH <- struct{}{}
	defer func() { <-__fieShuCH }()

	content := goo_utils.NewParams().Set("text", text)

	params := goo_utils.NewParams().
		Set("msg_type", "text").
		Set("content", content.Data())

	buf, err := goo_http_request.PostJson(hookUrl, params.JSON())
	if err != nil {
		return err
	}

	rst, _ := goo_utils.Byte(buf).Params()
	if msg := rst.Get("StatusMessage").String(); msg != "success" {
		return errors.New(msg)
	}

	return nil
}
