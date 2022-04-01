package goo_utils

import (
	"github.com/fvbock/endless"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func PProf(addr string, prefixOptions ...string) {
	prefix := "/goo/pprof"

	if l := len(prefixOptions); l > 0 {
		prefix = prefixOptions[0]
	}

	AsyncFunc(func() {
		engine := gin.Default()
		pprof.Register(engine, prefix)
		endless.NewServer(addr, engine).ListenAndServe()
	})
}
