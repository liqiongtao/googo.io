package goo_etcd

import (
	"fmt"
	goo_context "github.com/liqiongtao/googo.io/goo-context"
	"log"
	"testing"
	"time"
)

func TestInit(t *testing.T) {
	Init(Config{
		Endpoints: []string{"127.0.0.1:2379"},
		Username:  "root",
		Password:  "123456",
	})

	if _, err := Set("/goo/http-api/192.168.1.101:15001", "192.168.1.101:15001"); err != nil {
		log.Fatalln(err)
	}
	if _, err := Set("/goo/http-api/192.168.1.101:15002", "192.168.1.101:15002"); err != nil {
		log.Fatalln(err)
	}
	if _, err := SetTTL("/goo/http-api/192.168.1.101:15003", "192.168.1.101:15003", 3); err != nil {
		log.Fatalln(err)
	}

	for i := 0; i < 2; i++ {
		fmt.Println(GetString("/goo/http-api/"))
		fmt.Println(GetArray("/goo/http-api/"))
		fmt.Println(GetMap("/goo/http-api/"))
		time.Sleep(5 * time.Second)
	}

	Del("/goo/http-api/192.168.1.101:15002")
}

func TestRegisterService(t *testing.T) {
	Init(Config{
		Endpoints: []string{"127.0.0.1:2379"},
		Username:  "root",
		Password:  "123456",
	})

	err := RegisterService("/goo/http-api/node-1", "192.168.1.101:15002")
	fmt.Println(err)

	<-goo_context.Cancel().Done()
}

func TestWatch(t *testing.T) {
	Init(Config{
		Endpoints: []string{"127.0.0.1:2379"},
		Username:  "root",
		Password:  "123456",
	})

	go func() {
		for i := 0; i < 100; i++ {
			SetTTL(fmt.Sprintf("/goo/http-api/node-%d", i), fmt.Sprintf("192.168.1.%d", i), 5)
			time.Sleep(time.Second)
		}
	}()

	ch := Watch("/goo/http-api")
	for i := range ch {
		fmt.Println(i)
	}
}
