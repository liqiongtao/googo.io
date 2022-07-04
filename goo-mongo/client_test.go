package goo_mongo

import (
	"fmt"
	"testing"
)

func TestInit(t *testing.T) {
	conf := Config{
		Addr:     "122.228.113.230:27017",
		User:     "root",
		Password: "ab15eb8e12ea",
		Database: "test",
		AutoPing: false,
	}

	Init(conf)

	cli := GetClient()
	rst, err := cli.DB().Collection("user").InsertOne(cli.ctx, map[string]interface{}{
		"name": "hnatao",
		"age":  18,
	})
	fmt.Println(rst, err)
}
