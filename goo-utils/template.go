package goo_utils

import (
	"bytes"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"html/template"
)

func Template(text string, data interface{}) (string, error) {
	tpl, err := template.New("").Parse(text)
	if err != nil {
		goo_log.WithField("text", text).WithField("data", data).Error(err)
		return "", err
	}

	var b bytes.Buffer
	if err := tpl.Execute(&b, data); err != nil {
		goo_log.WithField("text", text).WithField("data", data).Error(err)
		return "", err
	}

	return b.String(), nil
}
