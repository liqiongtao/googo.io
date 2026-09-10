package main

import (
	"context"
	"time"

	goo_etcd "github.com/liqiongtao/googo.io/goo-etcd"
	goo_grpc "github.com/liqiongtao/googo.io/goo-grpc"
	pb_grpc_v1 "github.com/liqiongtao/googo.io/goo-grpc/test/proto"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
	"github.com/liqiongtao/googo.io/goocontext"
)

func main() {
	cli, err := goo_etcd.New(goo_etcd.Config{
		Endpoints: []string{"127.0.0.1:2379"},
	})
	if err != nil {
		panic(err)
	}

	s := goo_grpc.New(goo_grpc.Config{
		ServiceName:     "my-grpc",
		ServiceEndpoint: "my.grpc",
		Addr:            "127.0.0.1:10011",
	}).Register2Etcd(cli)

	pb_grpc_v1.RegisterGetterServer(s.Server, Server{})

	// 启动
	goo_utils.AsyncFunc(func() {
		s.Serve()
	})

	goocontext.Wait()
}

type Server struct {
	pb_grpc_v1.UnimplementedGetterServer
}

func (s Server) GetName(ctx context.Context, req *pb_grpc_v1.GetName_Request) (rsp *pb_grpc_v1.GetName_Response, err error) {
	rsp = &pb_grpc_v1.GetName_Response{}

	time.Sleep(time.Second)

	rsp.Name = time.Now().Format("15:04:05")

	return
}
