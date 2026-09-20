package goohttp

import (
	"encoding/json"
	"errors"
	goolog "github.com/liqiongtao/googo.io/goo-log"
	gooutils "github.com/liqiongtao/googo.io/goo-util"
	"time"
)

type Token struct {
	AppId     string `json:"app_id"`
	OpenId    int64  `json:"data"`
	NonceStr  string `json:"nonce"`
	Timestamp int64  `json:"ts"`
}

func (t *Token) Bytes() []byte {
	buf, _ := json.Marshal(t)
	return buf
}

func (t *Token) String() string {
	return string(t.Bytes())
}

func CreateToken(appId string, openid int64) (tokenStr string, err error) {
	token := &Token{
		AppId:     appId,
		OpenId:    openid,
		NonceStr:  gooutils.NonceStr(),
		Timestamp: time.Now().Unix(),
	}

	var (
		key    = gooutils.MD5([]byte(appId))
		iv     = key[8:24]
		encBuf []byte
	)

	encBuf, err = gooutils.AESCBCEncrypt(token.Bytes(), []byte(key), []byte(iv))
	if err != nil {
		goolog.Error(err.Error())
		return
	}

	tokenStr = gooutils.Base64Encode(encBuf)
	return
}

func ParseToken(tokenStr, appId string) (token *Token, err error) {
	var (
		tokenBuf = gooutils.Base64Decode(tokenStr)
		key      = gooutils.MD5([]byte(appId))
		iv       = key[8:24]
		b        []byte
	)

	b, err = gooutils.AESCBCDecrypt(tokenBuf, []byte(key), []byte(iv))
	if err != nil {
		goolog.Error(err.Error())
		return
	}

	token = new(Token)
	if err = json.Unmarshal(b, token); err != nil {
		goolog.Error(err.Error())
		return
	}
	if token.AppId != appId {
		err = errors.New("appid invalid")
	}
	return
}
