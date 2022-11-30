package goo_utils

import "encoding/json"

type M map[string]interface{}

func (m M) Json() []byte {
	b, _ := json.Marshal(&m)
	return b
}

func (m M) String() string {
	return string(m.Json())
}
