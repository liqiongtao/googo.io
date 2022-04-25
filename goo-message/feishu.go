package goo_message

import (
	"errors"
	goo_http_request "github.com/liqiongtao/googo.io/goo-http-request"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
	"runtime"
)

var FeiShu = &feishu{
	ch: make(chan struct{}, runtime.NumCPU()*2),
}

type feishu struct {
	ch chan struct{}
}

func (fs *feishu) Text(hookUrl string, text string) error {
	fs.ch <- struct{}{}
	defer func() { <-fs.ch }()

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
