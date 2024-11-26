package goo

import (
	"github.com/gin-gonic/gin"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
)

type Context interface {
	GinContext() *gin.Context
	SetGinContext(ctx *gin.Context)
	Log() *goo_log.Entry
}
