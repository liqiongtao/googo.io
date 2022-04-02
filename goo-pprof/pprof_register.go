package goo_pprof

import (
	"fmt"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	goo_context "github.com/liqiongtao/googo.io/goo-context"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
	"log"
	"net"
	"net/http"
)

// http://localhost:53404/goo/pprof
// go tool pprof http://localhost:53404/goo/pprof/heap
// go tool pprof http://localhost:53404/goo/pprof/profile
// go tool pprof http://localhost:53404/goo/pprof/block
// go tool pprof http://localhost:53404/goo/pprof/mutex
// go tool pprof http://localhost:53404/goo/pprof/goroutine
func Register(addr string, prefixArgs ...string) {
	prefix := "/goo/pprof"

	if l := len(prefixArgs); l > 0 {
		prefix = prefixArgs[0]
	}

	if addr == "" {
		addr = ":0"
	}

	var (
		listener net.Listener
	)

	goo_utils.AsyncFunc(func() {
		select {
		case <-goo_context.Cancel().Done():
			listener.Close()
		}
	})

	goo_utils.AsyncFunc(func() {
		var err error

		if listener, err = net.Listen("tcp", addr); err != nil {
			log.Println("Listen error:", err)
			return
		}
		defer listener.Close()

		gin.SetMode(gin.ReleaseMode)

		engine := gin.Default()

		pprof.Register(engine, prefix)

		l := listener.Addr().(*net.TCPAddr)
		log.Println(fmt.Sprintf("pprof http address=%s", l.String()))

		http.Serve(listener, engine)
	})
}
