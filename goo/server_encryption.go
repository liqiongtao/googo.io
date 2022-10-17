package goo

import (
	"encoding/hex"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
	"strings"
)

type Encryption struct {
	Key    string
	Secret string
}

func (enc *Encryption) Encode(b []byte) (str string, err error) {
	if l := len(b); l == 0 {
		return
	}
	var bts []byte
	bts, err = goo_utils.AESCBCEncrypt(b, []byte(enc.Key), []byte(enc.Secret))
	if err != nil {
		return
	}
	str = hex.EncodeToString(bts)
	return
}

func (enc *Encryption) Decode(str string) (b []byte, err error) {
	var bts []byte
	str = strings.ReplaceAll(str, "\"", "")
	bts, err = hex.DecodeString(str)
	if err != nil {
		return
	}
	b, err = goo_utils.AESCBCDecrypt(bts, []byte(enc.Key), []byte(enc.Secret))
	return
}
