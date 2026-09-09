package goo_http

import (
	"encoding/hex"
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
)

type Encryption struct {
	Key    string
	Secret string
}

func resolveEncryption(opts *options, c *gin.Context) (*Encryption, error) {
	if opts == nil || opts.encryptionFn == nil {
		return nil, errors.New("encryption not configured")
	}
	enc := opts.encryptionFn(c)
	if enc == nil {
		return nil, errors.New("encryption is nil")
	}
	return enc, nil
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
	str = strings.ReplaceAll(str, "\"", "")
	if l := len(str); l == 0 {
		return
	}
	var bts []byte
	bts, err = hex.DecodeString(str)
	if err != nil {
		return
	}
	b, err = goo_utils.AESCBCDecrypt(bts, []byte(enc.Key), []byte(enc.Secret))
	return
}
