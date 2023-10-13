package goo_utils

import (
	"bytes"
	"strings"
	"unicode"
)

// 多字符切割，默认支持逗号，分号，\n
func Split(s string, rs ...rune) []string {
	return strings.FieldsFunc(s, func(r rune) bool {
		if l := len(rs); l == 0 {
			return r == ',' || r == '，' || r == ';' || r == '；' || r == '\n'
		}
		for _, rr := range rs {
			if rr == r {
				return true
			}
		}
		return false
	})
}

// 驼峰转下划线
func Camel2Case(str string) string {
	var bf bytes.Buffer

	for i, r := range str {
		if !unicode.IsUpper(r) {
			bf.WriteRune(r)
			continue
		}
		if i > 0 {
			bf.WriteString("_")
		}
		bf.WriteRune(unicode.ToLower(r))
	}

	return bf.String()
}

// 下划线转驼峰
func Case2Camel(str string) string {
	str = strings.Replace(str, "_", " ", -1)
	str = strings.Title(str)
	return strings.Replace(str, " ", "", -1)
}
