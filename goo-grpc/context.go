package goo_grpc

import (
	"context"
	"github.com/gin-gonic/gin"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
	"google.golang.org/grpc/metadata"
	"strings"
	"time"
)

func Context(c *gin.Context) context.Context {
	ctx, _ := context.WithTimeout(context.Background(), 8*time.Second)
	md := metadata.New(map[string]string{})
	if c != nil {
		if v, ok := c.Get("__request_id"); ok {
			md.Set("x_trace_id", v.(string))
		}
		if v, ok := c.Get("__trace_id"); ok {
			md.Set("x_trace_id", v.(string))
		}
		if v, ok := c.Get("__server_name"); ok {
			md.Set("x_server_name", v.(string))
		}
		if v, ok := c.Get("__base_dir"); ok {
			arr := goo_utils.Trace(2, v.(string))
			if l := len(arr); l > 0 {
				md.Set("x_trace_info", strings.Join(arr, ", "))
			}
		}
	}
	return metadata.NewOutgoingContext(ctx, md)
}
