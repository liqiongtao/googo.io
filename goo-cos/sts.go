package goo_cos

import (
	"encoding/json"
	"fmt"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	goo_redis "github.com/liqiongtao/googo.io/goo-redis"
	sts "github.com/tencentyun/qcloud-cos-sts-sdk/go"
	"golang.org/x/sync/singleflight"
	"time"
)

// 获取临时密钥
// 文档: https://github.com/tencentyun/qcloud-cos-sts-sdk/tree/master/go
func GetCosSTSCredential(cfg StsConfig) (*sts.CredentialResult, error) {
	c := sts.NewClient(cfg.SecretId, cfg.SecretKey, nil)
	opt := &sts.CredentialOptions{
		Region:          cfg.Region,
		DurationSeconds: cfg.GetExpire(),
		Policy: &sts.CredentialPolicy{
			Statement: []sts.CredentialPolicyStatement{
				{
					Action:   cfg.GetAction(),
					Effect:   "allow",
					Resource: cfg.GetResource(),
				},
			},
		},
	}

	res, err := c.GetCredential(opt)
	if err != nil {
		goo_log.Error(err)
		return nil, err
	}
	if res.Error != nil {
		goo_log.Error(res.Error)
		return nil, res.Error
	}

	return res, nil
}

var (
	sfSts = singleflight.Group{}
)

// 获取临时密钥(带缓存)
func GetCosSTSCredentialWithCache(cfg StsConfig, redis *goo_redis.Client) (*sts.Credentials, error) {
	key := fmt.Sprintf("cos:sts:%s:%s:%s", cfg.Region, cfg.Appid, cfg.Bucket)

	result, err, _ := sfSts.Do(key, func() (interface{}, error) {
		var str string
		if redis != nil {
			str = redis.Get(key).Val()
		}
		if str != "" {
			// 从缓存中获取
			var credentials *sts.Credentials
			if err := json.Unmarshal([]byte(str), &credentials); err == nil {
				return credentials, nil
			}
		}

		// 从接口中获取
		res, err := GetCosSTSCredential(cfg)
		if err != nil {
			return nil, err
		}

		// 设置缓存
		if redis != nil && res.Credentials != nil {
			b, _ := json.Marshal(&res.Credentials)
			redis.Set(key, string(b), time.Duration(cfg.GetExpire())*time.Second)
		}

		return res.Credentials, nil
	})

	if err != nil {
		return nil, err
	}

	return result.(*sts.Credentials), nil
}
