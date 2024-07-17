package goo_utils

import (
	"bytes"
	"fmt"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"html/template"
	"strings"
)

func Template(text string, data interface{}) (string, []interface{}, error) {
	var args []interface{}
	var argsFunc = func(value interface{}) string {
		if v, ok := value.([]string); ok {
			var arr []string
			for _, vv := range v {
				arr = append(arr, "?")
				args = append(args, vv)
			}
			return strings.Join(arr, ",")
		}

		if v, ok := value.([]int64); ok {
			var arr []string
			for _, vv := range v {
				arr = append(arr, "?")
				args = append(args, vv)
			}
			return strings.Join(arr, ",")
		}

		args = append(args, value)
		return "?"
	}

	tpl := template.New("")
	tpl.Funcs(template.FuncMap{
		"Args": argsFunc,
		"LikeArgs": func(value interface{}) string {
			return argsFunc(fmt.Sprintf("%%%s%%", value))
		},
	})

	if _, err := tpl.Parse(text); err != nil {
		goo_log.WithField("text", text).WithField("data", data).Error(err)
		return "", []interface{}{}, err
	}

	var b bytes.Buffer
	if err := tpl.Execute(&b, data); err != nil {
		goo_log.WithField("text", text).WithField("data", data).Error(err)
		return "", []interface{}{}, err
	}

	return b.String(), args, nil
}
