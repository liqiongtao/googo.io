package goo_utils

import (
	"fmt"
	"testing"
)

func TestSplitStringArray(t *testing.T) {
	arr := []string{"a", "b", "c"}
	data := SplitStringArray(arr, 2)
	fmt.Println(data)
}

func TestSliceHas(t *testing.T) {
	arr := []string{"a", "b", "c"}
	has := SliceHas(arr, func(i int) bool {
		return arr[i] == "a"
	})
	fmt.Println(has)

	arr2 := []int64{1, 2, 3}
	has2 := SliceHas(arr, func(i int) bool {
		return arr2[i] == 5
	})
	fmt.Println(has2)
}

func TestSliceMap(t *testing.T) {
	arr := []string{"a", "b", "c"}
	SliceMap(arr, func(i int) {
		fmt.Println(arr[i])
	})
}

func TestSlice2UniqStrings(t *testing.T) {
	arr := []string{"a", "a", "b"}
	data := Slice2UniqStrings(arr, func(i int) string {
		return arr[i]
	})
	fmt.Println(data)
}

func TestSlice2UniqInt64s(t *testing.T) {
	arr := []int64{1, 2, 1}
	data := Slice2UniqInt64s(arr, func(i int) int64 {
		return arr[i]
	})
	fmt.Println(data)
}
