package goo_message

import (
	goo_http_request "googo.io/goo-http-request"
	goo_log "googo.io/goo-log"
	goo_utils "googo.io/goo-utils"
)

var FeiShu = new(feishu)

type feishu struct {
}

func (*feishu) Text(hookUrl string, text string) {
	content := goo_utils.Params{}.Set("text", text)

	params := goo_utils.Params{}.
		Set("msg_type", "text").
		Set("content", content.Data())

	buf, err := goo_http_request.PostJson(hookUrl, params.JSON())
	if err != nil {
		goo_log.Error(err.Error())
		return
	}
	
	rst, _ := goo_utils.Byte(buf).Params()
	if rst.Get("StatusMessage").String() != "success" {
		goo_log.Error(rst)
	}
}
