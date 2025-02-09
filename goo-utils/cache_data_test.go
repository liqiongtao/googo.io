package goo_utils

import (
	goo_context "github.com/liqiongtao/googo.io/goo-context"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"testing"
	"time"
)

func TestCacheDatafunc(t *testing.T) {
	for i := 0; i < 10; i++ {
		go cacheData1()
	}

	<-goo_context.WithCancel().Done()
}

func cacheData1() {
	key := "key-1"

	data, err := CacheDatafunc(key, func() (interface{}, error) {
		time.Sleep(time.Second)
		return time.Now().Format("15:04:05"), nil
	})

	if err != nil {
		goo_log.Error(err)
		return
	}

	goo_log.Info(data)
}
