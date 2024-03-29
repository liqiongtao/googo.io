package goo

import (
	"encoding/json"
	"errors"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
	"time"
)

type Token struct {
	AppId     string
	OpenId    int64
	NonceStr  string
	Timestamp int64
}

func (t *Token) Bytes() []byte {
	buf, _ := json.Marshal(t)
	return buf
}

func (t *Token) String() string {
	return string(t.Bytes())
}

func CreateToken(appId string, openId int64) (tokenStr string, err error) {
	token := &Token{
		AppId:     appId,
		OpenId:    openId,
		NonceStr:  goo_utils.NonceStr(),
		Timestamp: time.Now().Unix(),
	}

	var (
		key    = goo_utils.MD5([]byte(appId))
		iv     = key[8:24]
		encBuf []byte
	)

	encBuf, err = goo_utils.AESCBCEncrypt(token.Bytes(), []byte(key), []byte(iv))
	if err != nil {
		goo_log.Error(err.Error())
		return
	}

	tokenStr = goo_utils.Base64Encode(encBuf)
	return
}

func ParseToken(tokenStr, appId string) (token *Token, err error) {
	var (
		tokenBuf = goo_utils.Base64Decode(tokenStr)
		key      = goo_utils.MD5([]byte(appId))
		iv       = key[8:24]
		decBuf   []byte
	)

	decBuf, err = goo_utils.AESCBCDecrypt(tokenBuf, []byte(key), []byte(iv))
	if err != nil {
		goo_log.Error(err.Error())
		return
	}

	token = new(Token)
	if err = json.Unmarshal(decBuf, token); err != nil {
		goo_log.Error(err.Error())
		return
	}
	if token.AppId != appId {
		err = errors.New("appid invalid")
	}
	return
}
