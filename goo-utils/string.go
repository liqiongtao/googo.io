package goo_utils

import (
	"strings"
)

// 多字符切割，默认支持逗号，分号，\n
func Split(s string, rs ...rune) []string {
	return strings.FieldsFunc(s, func(r rune) bool {
		for _, rr := range rs {
			if rr == r {
				return true
			}
		}
		return r == ',' || r == '，' || r == ';' || r == '；' || r == '\n'
	})
}
