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
	res, err := GetCosSTSCredentialWithCache(cfg, redis)
	//fmt.Println(res, err)
	if err != nil {
		fmt.Println(err)
		return
	}

	c := CosClient(CosConfig{
		SecretId:     res.TmpSecretID,
		SecretKey:    res.TmpSecretKey,
		SessionToken: res.SessionToken,
		Bucket:       cfg.Bucket,
		Region:       cfg.Region,
	})

	if err := CosUpload("./1.log", "2025/1.log", c); err != nil {
		fmt.Println(err)
		return
	}

	b, err := CosGet("2025/1.log", c)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))

	if err := CosDelete("2025/1.log", c); err != nil {
		fmt.Println(err)
		return
	}
}
