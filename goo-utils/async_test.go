package goo_utils

import (
	"fmt"
	goo_context "github.com/liqiongtao/googo.io/goo-context"
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

	<-goo_context.WithCancel().Done()
}
