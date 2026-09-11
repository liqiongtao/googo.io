package goowechat

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	goolog "github.com/liqiongtao/googo.io/goo-log"
	goorequest "github.com/liqiongtao/googo.io/goo-request"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
)

type QrcodeParams struct {
	// 默认是主页，页面 page，例如 pages/index/index，根路径前不要填加 /，不能携带参数（参数请放在scene字段里），如果不填写这个字段，默认跳主页面。scancode_time为系统保留参数，不允许配置
	Page string `json:"page"`
	// 最大32个可见字符，只支持数字，大小写英文以及部分特殊字符：!#$&'()*+,/:;=?@-._~，其它字符请自行编码为合法字符（因不支持%，中文无法使用 urlencode 处理，请使用其他编码方式）
	Scene string `json:"scene"`
	// 默认430，二维码的宽度，单位 px，最小 280px，最大 1280px
	Width int `json:"width"`
	// 自动配置线条颜色，如果颜色依然是黑色，则说明不建议配置主色调，默认 false
	AutoColor bool `json:"auto_color"`
	// 默认是{"r":0,"g":0,"b":0} 。auto_color 为 false 时生效，使用 rgb 设置颜色 例如 {"r":"xxx","g":"xxx","b":"xxx"} 十进制表示
	LineColor goo_utils.M `json:"line_color"`
	// 默认是false，是否需要透明底色，为 true 时，生成透明底色的小程序
	IsHyaline bool `json:"is_hyaline"` // 是否透明
}

func (q *QrcodeParams) Json() []byte {
	if q.Width < 280 {
		q.Width = 600
	}
	if q.Width > 1280 {
		q.Width = 1280
	}

	b, _ := json.Marshal(q)
	return b
}

func Qrcode(appid, secret string, params QrcodeParams) (string, error) {
	accessToken := CGIToken(appid, secret).Get()

	urlstr := fmt.Sprintf(getwxacodeunlimit_url, accessToken)
	b, err := goorequest.PostJson(urlstr, params.Json())
	if err != nil {
		goolog.Error(err)
		return "", err
	}

	var errResp struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}
	if err = json.Unmarshal(b, &errResp); err == nil && errResp.ErrCode != 0 {
		goolog.Error(string(b))
		return "", errors.New(errResp.ErrMsg)
	}

	base64Str := base64.StdEncoding.EncodeToString(b)
	base64Image := fmt.Sprintf("data:image/png;base64,%s", base64Str)

	return base64Image, nil
}
