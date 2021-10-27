package main

import (
	"fmt"
	goo_kv "github.com/liqiongtao/googo.io/goo-kv"
)

var (
	endpoints = []string{
		"localhost:23791",
		"localhost:23792",
		"localhost:23793",
	}
)

func main() {
	goo_kv.SetAdapter(goo_kv.NewEtcdAdapter(endpoints))

	goo_kv.Set("/service/xzp/test-proj/test/node1", "127.0.0.1:13001", 5)
	goo_kv.Set("/service/xzp/test-proj/test/node2", "127.0.0.1:13002", 5)

	addr := goo_kv.Get("/service/xzp/test-proj/test/node1")
	fmt.Println(addr)

	addrs := goo_kv.GetMap("/service/xzp/test-proj/test")
	fmt.Println(addrs)

}
