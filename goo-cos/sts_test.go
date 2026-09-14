package goo_cos

import (
	"context"
	"fmt"
	"testing"

	gooredis "github.com/liqiongtao/googo.io/goo-redis"
)

func TestClient(t *testing.T) {
	redis, _ := gooredis.New(context.Background(), gooredis.Config{
		Addr:     "",
		Password: "",
		DB:       0,
		Prefix:   "",
	})

	cfg := StsConfig{
		SecretId:  "",
		SecretKey: "",
		Appid:     "",
		Bucket:    "",
		Region:    "",
	}
	res, err := STSCredentialWithCache(cfg, redis)
	fmt.Println(res, err)
}
