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
			md.Set("request_id", v.(string))
		}
		if v, ok := c.Get("__server_name"); ok {
			md.Set("server_name", v.(string))
		}
		if v, ok := c.Get("__base_dir"); ok {
			arr := goo_utils.Trace(2, v.(string))
			if l := len(arr); l > 0 {
				md.Set("request_trace", strings.Join(arr, ", "))
			}
		}
	}
	return metadata.NewOutgoingContext(ctx, md)
}
