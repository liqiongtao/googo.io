package goo_utils

import (
	"fmt"
	"testing"
)

func TestTemplate(t *testing.T) {
	var (
		tplUser = `
I like {{.Name}}

{{if gt .Age 30}} 
	to old 
{{else}} 
	{{.Age}} years old 
{{end}}

{{range $k, $v := .Likes}}
{{$k}} => {{$v}}
{{end}}

{{define "info"}}
	beijing
{{end}}

{{template "info"}}
`
	)

	m := M{
		"Name":  "hnatao",
		"Age":   30,
		"Likes": []string{"ping", "pang"},
	}

	fmt.Println(Template(tplUser, m))
}
