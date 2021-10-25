package goo_grpc

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
	"google.golang.org/grpc/metadata"
	"path"
	"strings"
	"time"
)

func Context(c *gin.Context) context.Context {
	ctx, _ := context.WithTimeout(context.Background(), 8*time.Second)
	md := metadata.New(map[string]string{})
	if c != nil {
		md.Set("request-id", fmt.Sprintf("%d", c.GetInt("__trace_id")))
		md.Set("request-server", c.GetString("__server_name"))
		files := goo_utils.Trace(2)
		if trimDir := path.Dir(c.GetString("__base_dir")); trimDir != "" {
			for i, f := range files {
				files[i] = strings.Replace(f, trimDir+"/", "", -1)
			}
		}
		md.Set("request-file", strings.Join(files, ", "))
	}
	return metadata.NewOutgoingContext(ctx, md)
}
