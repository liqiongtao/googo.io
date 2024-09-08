package main

import (
	"context"
	"fmt"
	goo_grpc "github.com/liqiongtao/googo.io/goo-grpc"
	pb_grpc_v1 "github.com/liqiongtao/googo.io/goo-grpc/my-proto"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"time"
)

func main() {
	cc, err := goo_grpc.Dial("127.0.0.1:10023")
	if err != nil {
		goo_log.Error(err)
		return
	}

	cli := pb_grpc_v1.NewGetterClient(cc)

	fmt.Println(time.Now().Format("15:04:05"))

	rsp, err := cli.GetName(context.TODO(), &pb_grpc_v1.GetName_Request{})
	fmt.Println(time.Now().Format("15:04:05"), rsp, err)
}
