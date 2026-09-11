package goo_utils

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	goo_log "github.com/liqiongtao/googo.io/goo-log"
)

func Template(text string, data any) (string, []any, error) {
	var args []any
	var argsFunc = func(value any) string {
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
		"args": argsFunc,
		"like": func(value any) string {
			return argsFunc(fmt.Sprintf("%%%s%%", value))
		},
	})

	if _, err := tpl.Parse(text); err != nil {
		goo_log.WithField("text", text).WithField("data", data).Error(err)
		return "", []any{}, err
	}

	var b bytes.Buffer
	if err := tpl.Execute(&b, data); err != nil {
		goo_log.WithField("text", text).WithField("data", data).Error(err)
		return "", []any{}, err
	}

	var lines []string
	for _, str := range strings.Split(b.String(), "\n") {
		if strings.TrimSpace(str) == "" {
			continue
		}
		lines = append(lines, str)
	}

	return strings.Join(lines, "\n"), args, nil
}
