package goo_utils

import (
	"fmt"
	"testing"
)

func TestId2Code(t *testing.T) {
	for i := 1; i < 1000; i++ {
		code := Id2Code(int64(i))
		id, _ := Code2Id(code)
		fmt.Println(id, code)
	}
}
