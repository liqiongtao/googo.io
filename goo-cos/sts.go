package goo_cos

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	goo_log "github.com/liqiongtao/googo.io/goo-log"
	goo_redis "github.com/liqiongtao/googo.io/goo-redis"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
	sts "github.com/tencentyun/qcloud-cos-sts-sdk/go"
	"golang.org/x/sync/singleflight"
)

// 获取临时密钥
// 文档: https://github.com/tencentyun/qcloud-cos-sts-sdk/tree/master/go
func STSCredential(cfg StsConfig) (*sts.CredentialResult, error) {
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

func stsCacheKey(cfg StsConfig) string {
	actions := append([]string{}, cfg.GetAction()...)
	sort.Strings(actions)
	fp := goo_utils.MD5([]byte(strings.Join(actions, ",") + "|" +
		fmt.Sprintf("%d", cfg.GetExpire()) + "|" +
		cfg.SecretId + "|" +
		strings.Join(cfg.GetResource(), ",")))
	return fmt.Sprintf("cos:sts:%s:%s:%s:%s", cfg.Region, cfg.Appid, cfg.Bucket, fp)
}

// 获取临时密钥(带缓存)
func STSCredentialWithCache(cfg StsConfig, redis *goo_redis.Client) (*sts.Credentials, error) {
	key := stsCacheKey(cfg)

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
		res, err := STSCredential(cfg)
		if err != nil {
			return nil, err
		}

		// 设置缓存（提前过期，避免临近失效仍被取用）
		if redis != nil && res.Credentials != nil {
			b, merr := json.Marshal(&res.Credentials)
			if merr != nil {
				return res.Credentials, nil
			}
			ttl := time.Duration(cfg.GetExpire()) * time.Second
			if ttl > 10*time.Minute {
				ttl -= 5 * time.Minute
			} else if ttl > time.Minute {
				ttl -= time.Minute
			}
			redis.Set(key, string(b), ttl)
		}

		return res.Credentials, nil
	})

	if err != nil {
		return nil, err
	}

	return result.(*sts.Credentials), nil
}
