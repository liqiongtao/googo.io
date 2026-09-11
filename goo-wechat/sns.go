package goowechat

import (
	"encoding/json"
	"errors"
	"fmt"
	goolog "github.com/liqiongtao/googo.io/goo-log"
	goorequest "github.com/liqiongtao/googo.io/goo-request"
	"net/url"
	"strings"
	"time"
)

// 获取H5授权链接
// authorizeUrl: 授权地址
// originUrl: 授权成功后的回跳地址
func Oauth2AuthorizeUrl(appid, authorizeUrl, originUrl, state string) string {
	var redirectUrl string
	if strings.Index(authorizeUrl, "?") != -1 {
		redirectUrl = url.QueryEscape(authorizeUrl + "&redirect_url=" + url.QueryEscape(originUrl))
	} else {
		redirectUrl = url.QueryEscape(authorizeUrl + "?redirect_url=" + url.QueryEscape(originUrl))
	}
	oauth2Url := fmt.Sprintf(oauth2_authorize_url, appid, redirectUrl, state)
	goolog.Debug(map[string]interface{}{
		"authorizeUrl": authorizeUrl,
		"originUrl":    originUrl,
		"redirectUrl":  redirectUrl,
		"oauth2Url":    oauth2Url,
	})
	return oauth2Url
}

type Oauth2AccessTokenResponse struct {
	AccessToken  string        `json:"access_token"`
	ExpiresIn    time.Duration `json:"expires_in"`
	RefreshToken string        `json:"refresh_token"`
	Openid       string        `json:"openid"`
	Unionid      string        `json:"unionid"`
	Scope        string        `json:"scope"`
	Errcode      int           `json:"errcode"`
	Errmsg       string        `json:"errmsg"`
}

func Oauth2AccessToken(appid, secret, code string) (*Oauth2AccessTokenResponse, error) {
	accessTokenUrl := fmt.Sprintf(sns_oauth2_accessToken_url, appid, secret, code)
	buf, err := goorequest.Get(accessTokenUrl)
	if err != nil {
		goolog.Error(err.Error())
		return nil, err
	}

	rsp := &Oauth2AccessTokenResponse{}
	if err := json.Unmarshal(buf, rsp); err != nil {
		goolog.Error(err.Error())
		return nil, err
	}
	if rsp.Errcode != 0 {
		goolog.Error(rsp.Errmsg)
		return nil, errors.New(rsp.Errmsg)
	}

	goolog.
		WithField("access_token", rsp.AccessToken).
		WithField("expire_in", rsp.ExpiresIn).
		WithField("openid", rsp.Openid).
		WithField("unionid", rsp.Unionid).
		Debug()

	return rsp, nil
}

type SnsUserInfoResponse struct {
	Openid   string `json:"openid"`
	Unionid  string `json:"unionid"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"headimgurl"`
	Gender   int    `json:"sex"`
	Country  string `json:"country"`
	Province string `json:"province"`
	City     string `json:"city"`
	Errcode  int    `json:"errcode"`
	Errmsg   string `json:"errmsg"`
}

func SnsUserInfo(accessToken, openid string) (*SnsUserInfoResponse, error) {
	userInfoUrl := fmt.Sprintf(sns_userinfo_url, accessToken, openid)
	buf, err := goorequest.Get(userInfoUrl)
	if err != nil {
		goolog.Error(err.Error())

		return nil, err
	}

	rsp := &SnsUserInfoResponse{}
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
		WithField("nickname", rsp.Nickname).
		WithField("gender", rsp.Gender).
		WithField("avatar", rsp.Avatar).
		WithField("province", rsp.Province).
		WithField("city", rsp.City).
		Debug()

	return rsp, nil
}
