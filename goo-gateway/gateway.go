package goo_gateway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/liqiongtao/googo.io/goo"
	pb_goo_v1 "github.com/liqiongtao/googo.io/goo-proto/v1"
	"google.golang.org/grpc/status"
)

type gateway struct {
	conf Config
}

func (g gateway) DoHandle(c *gin.Context) *goo.Response {
	service := c.Param("service")
	method := c.Param("method")

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, c.Request.Body); err != nil {
		return goo.Error(5000, "读取请求失败", err)
	}

	req := pb_goo_v1.Request{Data: buf.Bytes()}
	resp := pb_goo_v1.Response{}

	serviceName := fmt.Sprintf("/%s/%s/%s", g.conf.ServerName, service, g.conf.Env.Tag())
	cc := Client(serviceName)
	if cc == nil {
		return goo.Error(503, "服务不可用")
	}

	err := cc.Invoke(c, fmt.Sprintf("/%s/%s", service, method), &req, &resp)
	if s := status.Convert(err); s != nil {
		return goo.Error(s.Proto().Code, s.Proto().Message)
	}

	var data interface{}
	if l := len(resp.Data); l > 0 {
		if err = json.Unmarshal(resp.Data, &data); err != nil {
			return goo.Error(5001, "返回数据解析失败", err)
		}
	}

	return goo.Success(data)
}
