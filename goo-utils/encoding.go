package goo_utils

import (
	"bytes"
	"io"

	goolog "github.com/liqiongtao/googo.io/goo-log"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func GBK2UTF8(s string) string {
	r := bytes.NewReader([]byte(s))
	rr := transform.NewReader(r, simplifiedchinese.GBK.NewDecoder())
	buf, err := io.ReadAll(rr)
	if err != nil {
		goolog.WithField("str", s).Error(err.Error())
		return ""
	}
	return string(bytes.TrimSpace(buf))
}

func UTF82GBK(s string) string {
	r := bytes.NewReader([]byte(s))
	rr := transform.NewReader(r, simplifiedchinese.GBK.NewEncoder())
	buf, err := io.ReadAll(rr)
	if err != nil {
		goolog.WithField("str", s).Error(err.Error())
		return ""
	}
	return string(bytes.TrimSpace(buf))
}
