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

func TestTemplate2(t *testing.T) {
	m := M{
		"enterpriseId": 66,
		"name":         "123",
		"ids":          []int64{1, 2, 3},
	}

	sqlstr := `
SELECT * FROM u_user 
	WHERE enterprise_id = {{.enterpriseId|Args}} 
	and id IN ({{.ids|Args}})
	and name like {{.name|LikeArgs}}
`
	fmt.Println(Template(sqlstr, m))
}
