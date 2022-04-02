package goo_pprof

import (
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
	prefix := DEFAULT_PREFIX
	if l := len(prefixArgs); l > 0 {
		prefix = prefixArgs[0]
	}

	// 随机端口
	if addr == "" {
		addr = ":0"
	}

	var (
		listener net.Listener
	)

	// 监听信号，kill -1 时，关闭链接
	goo_utils.AsyncFunc(func() {
		select {
		case <-goo_context.Cancel().Done():
			listener.Close()
		}
	})

	goo_utils.AsyncFunc(func() {
		var err error

		// 启动监听
		if listener, err = net.Listen("tcp", addr); err != nil {
			log.Println("Listen error:", err)
			return
		}
		defer listener.Close()

		gin.SetMode(gin.ReleaseMode)

		engine := gin.Default()

		// 注册
		pprof.Register(engine, prefix)

		// 启动服务
		http.Serve(listener, engine)
	})
}
