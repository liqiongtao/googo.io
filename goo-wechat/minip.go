package goowechat

import (
	"encoding/json"
	"errors"
	"fmt"
	goolog "github.com/liqiongtao/googo.io/goo-log"
	goorequest "github.com/liqiongtao/googo.io/goo-request"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
)

// ---------------------------------
// -- 小程序登录
// ---------------------------------

type JsCode2SessionResponse struct {
	Openid     string `json:"openid"`
	Unionid    string `json:"unionid"`
	SessionKey string `json:"session_key"`
	Errcode    int    `json:"errcode"`
	Errmsg     string `json:"errmsg"`
}

func JsCode2Session(appid, secret, code string) (*JsCode2SessionResponse, error) {
	jscode2sess_url := fmt.Sprintf(sns_jsscode2sess_url, appid, secret, code)
	buf, err := goorequest.Get(jscode2sess_url)
	if err != nil {
		goolog.Error(err.Error())
		return nil, err
	}

	rsp := &JsCode2SessionResponse{}
	if err := json.Unmarshal(buf, rsp); err != nil {
		goolog.Error(err.Error())
		return nil, err
	}
	if rsp.Errcode != 0 {
		goolog.Error(rsp.Errmsg)
		return nil, errors.New(rsp.Errmsg)
	}

	goolog.
		WithField("openid", rsp.Openid).
		WithField("unionid", rsp.Unionid).
		WithField("session_key", rsp.SessionKey).
		Debug()

	return rsp, nil
}

// ---------------------------------
// -- 解析用户数据
// ---------------------------------

type MinipUserInfoResponse struct {
	Openid   string `json:"openid"`
	Unionid  string `json:"unionid"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatarUrl"`
	Gender   int    `json:"gender"`
	Country  string `json:"country"`
	Province string `json:"province"`
	City     string `json:"city"`
}

func MinipUserInfo(sessionKey, encryptedData, iv string) (*MinipUserInfoResponse, error) {
	data := goo_utils.Base64Decode(encryptedData)
	key := goo_utils.Base64Decode(sessionKey)

	buf, err := goo_utils.AESCBCDecrypt(data, key, goo_utils.Base64Decode(iv))
	if err != nil {
		goolog.Error(err.Error())
		return nil, err
	}

	rsp := &MinipUserInfoResponse{}
	if err = json.Unmarshal(buf, rsp); err != nil {
		goolog.Error(err.Error())
		return nil, err
	}

	goolog.
		WithField("openid", rsp.Openid).
		WithField("unionid", rsp.Unionid).
		WithField("nickname", rsp.Nickname).
		WithField("gender", rsp.Gender).
		WithField("avatar", rsp.Avatar).
		WithField("province", rsp.Province).
		WithField("city", rsp.City).
		Debug()

	return rsp, nil
}

// ---------------------------------
// -- 解析手机号
// ---------------------------------

type WXMobileData struct {
	PhoneNumber     string `json:"phoneNumber"`
	PurePhoneNumber string `json:"purePhoneNumber"`
	CountryCode     string `json:"countryCode"`
	Watermark       struct {
		Appid     string `json:"appid"`
		Timestamp int64  `json:"timestamp"`
	} `json:"watermark"`
}

func MinipMobile(sessionKey, encryptedData, iv string) (*WXMobileData, error) {
	// 解析数据
	encryptedDataRaw := goo_utils.Base64Decode(encryptedData)
	ivRaw := goo_utils.Base64Decode(iv)
	key := goo_utils.Base64Decode(sessionKey)
	buf, err := goo_utils.AESCBCDecrypt(encryptedDataRaw, key, ivRaw)
	if err != nil {
		return nil, err
	}
	// 获取手机号
	dt := &WXMobileData{}
	if err = json.Unmarshal(buf, dt); err != nil {
		return nil, err
	}
	return dt, nil
}

// ---------------------------------
// -- 发送订阅消息
// ---------------------------------

func SendTemplateMessage(appid, secret, openid, templateId, page string, data interface{}) error {
	accessToken := CGIToken(appid, secret).Get()
	params := goo_utils.M{
		"access_token": accessToken,
		"template_id":  templateId,
		"page":         page,
		"touser":       openid,
		"data":         data,
	}

	messageTplSendUrl := fmt.Sprintf(message_tpl_send_url, accessToken)
	b, err := goorequest.PostJson(messageTplSendUrl, params.Json())
	if err != nil {
		goolog.Error(err.Error())
		return err
	}

	l := goolog.WithField("params", params.String()).WithField("result", string(b))

	p, err := goo_utils.Byte(b).Params()
	if err != nil {
		l.Error("发送订阅消息", err)
		return err
	}

	if p.Get("errcode").Int() != 0 {
		l.Error("发送订阅消息")
		return errors.New(p.Get("errmsg").String())
	}

	l.Debug("发送订阅消息")

	return nil
}
