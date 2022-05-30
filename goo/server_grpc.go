package goo

import (
	"context"
	"github.com/gin-gonic/gin"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
	"google.golang.org/grpc/metadata"
	"strings"
)

func GrpcContext(c *gin.Context) context.Context {
	md := metadata.New(map[string]string{})
	if c != nil {
		if v := c.GetString("__server_name"); v != "" {
			md.Set("server-name", v)
		}
		if v := c.GetString("__env"); v != "" {
			md.Set("env", v)
		}
		if v := requestId(c); v != "" {
			md.Set("trace-id", v)
		}
		if v := requestId(c); v != "" {
			arr := goo_utils.Trace(2)
			if l := len(arr); l > 0 {
				md.Set("caller", strings.Join(arr, ", "))
			}
		}
	}
	return metadata.NewOutgoingContext(context.TODO(), md)
}
