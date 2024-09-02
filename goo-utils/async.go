package goo_utils

import (
	"context"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"runtime"
	"sync"
	"time"
)

// 捕获panic
func Recovery() {
	if r := recover(); r != nil {
		goo_log.Error(r)
	}
}

// 异步执行（安全）
func AsyncFunc(fn func()) {
	go func() {
		defer Recovery()
		fn()
	}()
}

// 异步执行（安全）
func AsyncFuncWithTimeout(fn func(), d time.Duration) {
	ctx, cancel := context.WithCancel(context.TODO())

	go func() {
		defer Recovery()
		defer func() { cancel() }()
		fn()
	}()

	select {
	case <-ctx.Done():
	case <-time.Tick(d):
	}
}

// 异步并发执行（安全）
func AsyncFuncGroup(fns ...func()) {
	var (
		wg sync.WaitGroup
		ch = make(chan struct{}, runtime.NumCPU())
	)

	for _, fn := range fns {
		wg.Add(1)
		ch <- struct{}{}

		func(fn func()) {
			AsyncFunc(func() {
				defer wg.Done()
				defer func() { <-ch }()
				fn()
			})
		}(fn)
	}

	wg.Wait()
}
