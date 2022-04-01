package goo_utils

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func PProf(addr string, prefixOptions ...string) {
	prefix := "/goo/pprof"

	if l := len(prefixOptions); l > 0 {
		prefix = prefixOptions[0]
	}

	AsyncFunc(func() {
		r := gin.Default()
		pprof.Register(r, prefix)
		r.Run(addr)
	})
}
