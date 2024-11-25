package main

import (
	"context"
	"time"

	goo_context "github.com/liqiongtao/googo.io/goo-context"
	goo_etcd "github.com/liqiongtao/googo.io/goo-etcd"
	goo_grpc "github.com/liqiongtao/googo.io/goo-grpc"
	pb_grpc_v1 "github.com/liqiongtao/googo.io/goo-grpc/test/proto"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func main() {
	cli, err := goo_etcd.New(goo_etcd.Config{
		Endpoints: []string{"127.0.0.1:2379"},
	})
	if err != nil {
		panic(err)
	}

	cc, err := goo_grpc.DialContextWithEtcd(context.TODO(), "my-grpc", cli)
	if err != nil {
		goo_log.Error(err)
		return
	}

	c := pb_grpc_v1.NewGetterClient(cc)

	goo_utils.AsyncFunc(func() {
		uuid := goo_utils.UUID()
		for {
			ctx := metadata.NewOutgoingContext(context.TODO(), metadata.New(map[string]string{"trace-id": uuid}))

			rsp, err := c.GetName(ctx, &pb_grpc_v1.GetName_Request{Name: "hnatao"})
			if s := status.Convert(err); s.Code() != codes.OK {
				goo_log.Error("request error", s.Message())
				time.Sleep(time.Second)
				continue
			}

			goo_log.Info(rsp.Name)
		}
	})

	<-goo_context.WithCancel().Done()
}
