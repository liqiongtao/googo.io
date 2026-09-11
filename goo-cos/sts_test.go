package goo_cos

import (
	"fmt"
	gooredis "github.com/liqiongtao/googo.io/goo-redis"
	"testing"
)

func TestClient(t *testing.T) {
	redis, _ := gooredis.New(gooredis.Config{
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
