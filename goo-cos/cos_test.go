package goo_cos

import (
	"fmt"
	goo_redis "github.com/liqiongtao/googo.io/goo-redis"
	"testing"
)

func TestCosClient(t *testing.T) {
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
	//fmt.Println(res, err)
	if err != nil {
		fmt.Println(err)
		return
	}

	c := NewCosClient(CosConfig{
		SecretId:     res.TmpSecretID,
		SecretKey:    res.TmpSecretKey,
		SessionToken: res.SessionToken,
		Bucket:       cfg.Bucket,
		Region:       cfg.Region,
	})

	if err := c.Upload("./1.log", "2025/1.log"); err != nil {
		fmt.Println(err)
		return
	}

	b, err := c.Get("2025/1.log")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))

	if err := c.Delete("2025/1.log"); err != nil {
		fmt.Println(err)
		return
	}
}
