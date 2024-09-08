package main

import (
	"context"
	"fmt"
	goo_grpc "github.com/liqiongtao/googo.io/goo-grpc"
	pb_grpc_v1 "github.com/liqiongtao/googo.io/goo-grpc/my-proto"
	"time"
)

func main() {
	s := goo_grpc.New(goo_grpc.Config{
		ServiceName:     "my-grpc",
		ServiceEndpoint: "my.grpc",
		Addr:            "127.0.0.1:10023",
	})

	pb_grpc_v1.RegisterGetterServer(s.Server, Server{})

	s.Serve()
}

type Server struct {
}

func (s Server) GetName(ctx context.Context, request *pb_grpc_v1.GetName_Request) (*pb_grpc_v1.GetName_Response, error) {
	fmt.Println("-------0---------")
	time.Sleep(5 * time.Minute)

	return &pb_grpc_v1.GetName_Response{}, nil
}
