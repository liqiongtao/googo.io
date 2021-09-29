package goo_grpc

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
	"path"
	"runtime"
	"strings"
	"time"
)

func Context(c *gin.Context) context.Context {
	ctx, _ := context.WithTimeout(context.Background(), 8*time.Second)
	md := metadata.New(map[string]string{})
	if c != nil {
		md.Set("request-id", fmt.Sprintf("%d", c.GetInt("__trace_id")))
		md.Set("request-server", c.GetString("__server_name"))
		_, file, line, _ := runtime.Caller(1)
		file = strings.Replace(file, path.Dir(c.GetString("__base_dir"))+"/", "", -1)
		md.Set("request-file", fmt.Sprintf("%s %dL", file, line))
	}
	return metadata.NewOutgoingContext(ctx, md)
}
