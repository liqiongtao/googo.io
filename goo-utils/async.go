package goo_utils

import (
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"sync"
)

// 捕获panic
func Recovery() {
	if err := recover(); err != nil {
		goo_log.WithTrace().Error(err)
	}
}

// 异步执行（安全）
func AsyncFunc(fn func()) {
	go func() {
		defer Recovery()
		fn()
	}()
}

// 异步并发执行（安全）
func AsyncFuncGroup(fns ...func()) {
	var wg sync.WaitGroup

	for _, fn := range fns {
		wg.Add(1)
		func(fn func()) {
			AsyncFunc(func() {
				defer wg.Done()
				fn()
			})
		}(fn)
	}

	wg.Wait()
}
