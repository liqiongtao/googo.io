package goo_utils

import (
	"fmt"
	"testing"
)

func TestJson2Params(t *testing.T) {
	str := `{"name":"hnatao"}`

	arr := []Params{}

	for i := 0; i < 3; i++ {
		p, _ := Json2Params([]byte(str))
		arr = append(arr, p)
	}

	arr[1].Set("info", map[string]interface{}{"sex": "男"})

	fmt.Println(arr)
	fmt.Println(arr[1])
	fmt.Println(arr[1].Get("name"))
	fmt.Println(arr[1].Get("info.sex"))
}
