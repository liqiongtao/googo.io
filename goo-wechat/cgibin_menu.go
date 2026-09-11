package goowechat

import (
	"encoding/json"
	"errors"
	"fmt"
	goolog "github.com/liqiongtao/googo.io/goo-log"
	goorequest "github.com/liqiongtao/googo.io/goo-request"
)

func MenuCreate(appid, secret, content string) error {
	accessToken := CGIToken(appid, secret).Get()

	menuCreateUrl := fmt.Sprintf(menu_create_url, accessToken)
	buf, err := goorequest.PostJson(menuCreateUrl, []byte(content))
	if err != nil {
		goolog.Error(err.Error())
		return err
	}

	rst := struct {
		ErrorCode int    `json:"errorcode"`
		ErrMsg    string `json:"errmsg"`
	}{}
	if err := json.Unmarshal(buf, &rst); err != nil {
		goolog.Error(err.Error())
		return err
	}
	if rst.ErrorCode != 0 {
		goolog.Error(rst.ErrMsg)
		return errors.New(rst.ErrMsg)
	}

	return nil
}

func MenuGet(appid, secret string) (string, error) {
	accessToken := CGIToken(appid, secret).Get()

	menuGetrl := fmt.Sprintf(menu_get_url, accessToken)
	buf, err := goorequest.Get(menuGetrl)
	if err != nil {
		goolog.Error(err.Error())
		return "", err
	}

	rst := struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}{}
	if err := json.Unmarshal(buf, &rst); err != nil {
		goolog.Error(err.Error())
		return "", err
	}
	if rst.ErrCode != 0 {
		goolog.Error(rst.ErrMsg)
		return "", errors.New(rst.ErrMsg)
	}

	return string(buf), nil
}

func MenuDelete(appid, secret string) error {
	accessToken := CGIToken(appid, secret).Get()

	menuDeleteUrl := fmt.Sprintf(menu_del_url, accessToken)
	buf, err := goorequest.PostJson(menuDeleteUrl, nil)
	if err != nil {
		goolog.Error(err.Error())
		return err
	}

	rst := struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}{}
	if err := json.Unmarshal(buf, &rst); err != nil {
		goolog.Error(err.Error())
		return err
	}
	if rst.ErrCode != 0 {
		goolog.Error(rst.ErrMsg)
		return errors.New(rst.ErrMsg)
	}
	return nil
}
