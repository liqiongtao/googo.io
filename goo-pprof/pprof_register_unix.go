package goo_pprof

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	goo_context "github.com/liqiongtao/googo.io/goo-context"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
	"log"
	"net"
	"net/http"
	"os"
)

// curl --unix-socket pprof.sock http://localhost/goo/pprof/heap -o heap.pprof.gz
// go tool pprof heap.pprof.gz

// curl --unix-socket pprof.sock http://localhost/goo/pprof/profile -o profile.pprof.gz
// go tool pprof profile.pprof.gz

// curl --unix-socket pprof.sock http://localhost/goo/pprof/block -o block.pprof.gz
// go tool pprof block.pprof.gz

// curl --unix-socket pprof.sock http://localhost/goo/pprof/mutex -o mutex.pprof.gz
// go tool pprof mutex.pprof.gz

// curl --unix-socket pprof.sock http://localhost/goo/pprof/goroutine -o goroutine.pprof.gz
// go tool pprof goroutine.pprof.gz

func RegisterUnix(prefixArgs ...string) {
	prefix := "/goo/pprof"

	if l := len(prefixArgs); l > 0 {
		prefix = prefixArgs[0]
	}

	var (
		file     = "pprof.sock"
		listener net.Listener
	)

	if _, err := os.Stat(file); err == nil {
		os.Remove(file)
	}

	goo_utils.AsyncFunc(func() {
		select {
		case <-goo_context.Cancel().Done():
			listener.Close()
			os.Remove(file)
		}
	})

	goo_utils.AsyncFunc(func() {
		var err error

		if listener, err = net.Listen("unix", file); err != nil {
			log.Println("Listen error:", err)
			return
		}
		defer listener.Close()
		defer os.Remove(file)

		gin.SetMode(gin.ReleaseMode)

		engine := gin.Default()

		pprof.Register(engine, prefix)

		http.Serve(listener, engine)
	})
}
