package goo_utils

import "reflect"

func SplitStringArray(arr []string, size int) (list [][]string) {
	l := len(arr)

	if l == 0 {
		list = make([][]string, 0)
		return
	}

	if l < size {
		list = [][]string{arr}
		return
	}

	var (
		offset int
	)

	for {
		if offset+size >= l {
			list = append(list, arr[offset:])
			break
		}

		list = append(list, arr[offset:offset+size])

		offset += size
	}

	return
}

func SplitIntArray(arr []int, size int) (list [][]int) {
	l := len(arr)

	if l == 0 {
		list = make([][]int, 0)
		return
	}

	if l < size {
		list = [][]int{arr}
		return
	}

	var (
		offset int
	)

	for {
		if offset+size >= l {
			list = append(list, arr[offset:])
			break
		}

		list = append(list, arr[offset:offset+size])

		offset += size
	}

	return
}

func SplitInt64Array(arr []int64, size int) (list [][]int64) {
	l := len(arr)

	if l == 0 {
		list = make([][]int64, 0)
		return
	}

	if l < size {
		list = [][]int64{arr}
		return
	}

	var (
		offset int
	)

	for {
		if offset+size >= l {
			list = append(list, arr[offset:])
			break
		}

		list = append(list, arr[offset:offset+size])

		offset += size
	}

	return
}

func SplitArray(arr []interface{}, size int) (list [][]interface{}) {
	l := len(arr)

	if l == 0 {
		list = make([][]interface{}, 0)
		return
	}

	if l < size {
		list = [][]interface{}{arr}
		return
	}

	var (
		offset int
	)

	for {
		if offset+size >= l {
			list = append(list, arr[offset:])
			break
		}

		list = append(list, arr[offset:offset+size])

		offset += size
	}

	return
}

func SliceHas(x any, f func(i int) bool) bool {
	rv := reflect.ValueOf(x)
	for i := 0; i < rv.Len(); i++ {
		if f(i) {
			return true
		}
	}
	return false
}

func SliceMap(x any, f func(i int)) {
	rv := reflect.ValueOf(x)
	for i := 0; i < rv.Len(); i++ {
		f(i)
	}
	return
}

func Slice2UniqStrings(x any, f func(i int) string) (data []string) {
	data = []string{}
	m := map[string]struct{}{}
	rv := reflect.ValueOf(x)
	for i := 0; i < rv.Len(); i++ {
		v := f(i)
		if _, ok := m[v]; !ok {
			m[v] = struct{}{}
			data = append(data, v)
		}
	}
	return
}

func Slice2UniqInt64s(x any, f func(i int) int64) (data []int64) {
	data = []int64{}
	m := map[int64]struct{}{}
	rv := reflect.ValueOf(x)
	for i := 0; i < rv.Len(); i++ {
		v := f(i)
		if _, ok := m[v]; !ok {
			m[v] = struct{}{}
			data = append(data, v)
		}
	}
	return
}

func Slice2UniqInt32s(x any, f func(i int) int32) (data []int32) {
	data = []int32{}
	m := map[int32]struct{}{}
	rv := reflect.ValueOf(x)
	for i := 0; i < rv.Len(); i++ {
		v := f(i)
		if _, ok := m[v]; !ok {
			m[v] = struct{}{}
			data = append(data, v)
		}
	}
	return
}

func Slice2UniqInts(x any, f func(i int) int) (data []int) {
	data = []int{}
	m := map[int]struct{}{}
	rv := reflect.ValueOf(x)
	for i := 0; i < rv.Len(); i++ {
		v := f(i)
		if _, ok := m[v]; !ok {
			m[v] = struct{}{}
			data = append(data, v)
		}
	}
	return
}

func Slice2UniqFloat64s(x any, f func(i int) float64) (data []float64) {
	data = []float64{}
	m := map[float64]struct{}{}
	rv := reflect.ValueOf(x)
	for i := 0; i < rv.Len(); i++ {
		v := f(i)
		if _, ok := m[v]; !ok {
			m[v] = struct{}{}
			data = append(data, v)
		}
	}
	return
}

func Slice2UniqFloat32s(x any, f func(i int) float32) (data []float32) {
	data = []float32{}
	m := map[float32]struct{}{}
	rv := reflect.ValueOf(x)
	for i := 0; i < rv.Len(); i++ {
		v := f(i)
		if _, ok := m[v]; !ok {
			m[v] = struct{}{}
			data = append(data, v)
		}
	}
	return
}
