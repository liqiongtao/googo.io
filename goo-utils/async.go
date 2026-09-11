package goo_utils

import (
	"context"
	"runtime"
	"sync"
	"time"

	goolog "github.com/liqiongtao/googo.io/goo-log"
)

// 捕获panic
func Recovery() {
	if r := recover(); r != nil {
		goolog.Error(r)
	}
}

// 异步执行（安全）
func AsyncFunc(fn func()) {
	go func() {
		defer Recovery()
		fn()
	}()
}

// AsyncFuncWithTimeout 在超时时间内等待 fn 完成。
// 超时后取消 ctx，fn 应监听 ctx.Done() 以尽快退出；超时返回后 fn 可能仍在运行直至其响应取消。
func AsyncFuncWithTimeout(fn func(ctx context.Context), d time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), d)
	defer cancel()

	done := make(chan struct{})
	go func() {
		defer Recovery()
		defer close(done)
		fn(ctx)
	}()

	select {
	case <-done:
	case <-ctx.Done():
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
