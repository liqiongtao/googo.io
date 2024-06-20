package goo_utils

import (
	"fmt"
	"testing"
)

func TestStr2Time(t *testing.T) {
	fmt.Println(Str2Time("2024-06-01 10:01:02"))
	fmt.Println(Str2Time("2024-6-01 10:1:02"))
	fmt.Println(Str2Time("2024-6-1 10:1:2"))
	fmt.Println(Str2Time("2024-6-01 10:1"))
	fmt.Println(Str2Time("2024-6-1 10:1"))

	fmt.Println(Str2Time("2024/06/01 10:01:02"))
	fmt.Println(Str2Time("2024/6/01 10:1:02"))
	fmt.Println(Str2Time("2024/6/1 10:1:2"))
	fmt.Println(Str2Time("2024/6/01 10:1"))
	fmt.Println(Str2Time("2024/6/1 10:1"))

	fmt.Println(Str2Time("20240601100102"))
	fmt.Println(Str2Time("20240601_100102"))

	fmt.Println(Str2Time("2024-6-1"))
	fmt.Println(Str2Time("2024/6/1"))

	fmt.Println(Str2Time("10:01:02"))
	fmt.Println(Str2Time("10:1:02"))
	fmt.Println(Str2Time("10:1:2"))

	fmt.Println(Str2Time("10:01"))
	fmt.Println(Str2Time("10:1"))
	fmt.Println(Str2Time("10:1"))
}
