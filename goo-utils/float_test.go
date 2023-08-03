package goo_utils

import (
	"fmt"
	"testing"
)

func TestFloat64_ToFixed(t *testing.T) {
	fmt.Println(Float64(0.123456789).ToFixed(0))
	fmt.Println(Float64(0.123456789).ToFixed(1))
	fmt.Println(Float64(0.123456789).ToFixed(2))
	fmt.Println(Float64(0.123456789).ToFixed(3))
	fmt.Println(Float64(0.123456789).ToFixed(4))
	fmt.Println(Float64(0.123456789).ToFixed(5))

	fmt.Println(Float64(0.123456789).ToPercent(0))
	fmt.Println(Float64(0.123456789).ToPercent(1))
	fmt.Println(Float64(0.123456789).ToPercent(2))
	fmt.Println(Float64(0.123456789).ToPercent(3))
	fmt.Println(Float64(0.123456789).ToPercent(4))
	fmt.Println(Float64(0.123456789).ToPercent(5))
}
