package goo_utils

import (
	"fmt"
	"testing"
	"time"
)

func TestAsyncFuncWithTimeout(t *testing.T) {
	AsyncFuncWithTimeout(func() {
		for i := 0; i < 5; i++ {
			fmt.Println(i)
			time.Sleep(time.Second)
		}
	}, 3*time.Second)

	fmt.Println("done")
}
