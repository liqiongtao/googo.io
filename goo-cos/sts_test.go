package goo_cos

import (
	"fmt"
	goo_redis "github.com/liqiongtao/googo.io/goo-redis"
	"testing"
)

func TestClient(t *testing.T) {
	redis, _ := goo_redis.New(goo_redis.Config{
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
