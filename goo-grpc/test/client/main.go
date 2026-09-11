package main

import (
	"context"
	"time"

	"github.com/liqiongtao/googo.io/goo-context"
	gooetcd "github.com/liqiongtao/googo.io/goo-etcd"
	googrpc "github.com/liqiongtao/googo.io/goo-grpc"
	pb_grpc_v1 "github.com/liqiongtao/googo.io/goo-grpc/test/proto"
	goolog "github.com/liqiongtao/googo.io/goo-log"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func main() {
	cli, err := gooetcd.New(gooetcd.Config{
		Endpoints: []string{"127.0.0.1:2379"},
	})
	if err != nil {
		panic(err)
	}

	cc, err := googrpc.DialContextWithEtcd(context.TODO(), "my-grpc", cli)
	if err != nil {
		goolog.Error(err)
		return
	}

	c := pb_grpc_v1.NewGetterClient(cc)

	goo_utils.AsyncFunc(func() {
		uuid := goo_utils.UUID()
		for {
			ctx := metadata.NewOutgoingContext(context.TODO(), metadata.New(map[string]string{"trace-id": uuid}))

			rsp, err := c.GetName(ctx, &pb_grpc_v1.GetName_Request{Name: "hnatao"})
			if s := status.Convert(err); s.Code() != codes.OK {
				goolog.Error("request error", s.Message())
				time.Sleep(time.Second)
				continue
			}

			goolog.Info(rsp.Name)
		}
	})

	goocontext.Wait()
}
