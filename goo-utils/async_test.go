package goo_utils

import (
	"fmt"
	"github.com/liqiongtao/googo.io/goocontext"
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

	<-goocontext.Root().Done()
}
