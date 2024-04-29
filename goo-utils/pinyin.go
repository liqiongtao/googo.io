package goo_utils

import (
	"github.com/mozillazg/go-pinyin"
	"strings"
)

func Pinyin(str string) []string {
	return pinyin.LazyPinyin(str, pinyin.NewArgs())
}

func PinyinStr(str string) string {
	return strings.Join(Pinyin(str), "")
}
