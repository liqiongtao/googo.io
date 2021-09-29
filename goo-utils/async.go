package goo_utils

func AsyncFunc(fn func()) {
	go func(fn func()) {
		defer Recover()
		fn()
	}(fn)
}
