package goo_utils

import "testing"

type User struct {
	Name string `json:"name"`
}

func TestStr2Struct(t *testing.T) {
	str := `{"name": "hnatao"}`
	user, err := Str2Struct[User](str)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(user.Name)
}
