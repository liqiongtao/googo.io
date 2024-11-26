package goo

import (
	"context"
	"github.com/gin-gonic/gin"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
)

type Context struct {
	*gin.Context
}

func (ctx Context) Log() *goo_log.Entry {
	return goo_log.WithField("trace_id", RequestId(ctx.Context))
}

func (ctx Context) WithTraceId() context.Context {
	return context.WithValue(context.TODO(), "trace_id", RequestId(ctx.Context))
}
