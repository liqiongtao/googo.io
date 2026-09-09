package goo_utils

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestAsyncFuncWithTimeout(t *testing.T) {
	AsyncFuncWithTimeout(func(ctx context.Context) {
		for i := 0; i < 5; i++ {
			select {
			case <-ctx.Done():
				fmt.Println("cancelled")
				return
			default:
			}
			fmt.Println(i)
			time.Sleep(time.Second)
		}
	}, 3*time.Second)

	fmt.Println("done")
}
