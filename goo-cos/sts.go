package goo_cos

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	goolog "github.com/liqiongtao/googo.io/goo-log"
	gooredis "github.com/liqiongtao/googo.io/goo-redis"
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
		goolog.Error(err)
		return nil, err
	}
	if res.Error != nil {
		goolog.Error(res.Error)
		return nil, res.Error
	}
	if res.Credentials == nil {
		err = fmt.Errorf("sts credentials is nil")
		goolog.Error(err)
		return nil, err
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
func STSCredentialWithCache(cfg StsConfig, redis *gooredis.Client) (*sts.Credentials, error) {
	key := stsCacheKey(cfg)

	result, err, _ := sfSts.Do(key, func() (any, error) {
		if redis != nil {
			cmd := redis.Get(key)
			if cerr := cmd.Err(); cerr == nil {
				var credentials *sts.Credentials
				if uerr := json.Unmarshal([]byte(cmd.Val()), &credentials); uerr == nil && credentials != nil {
					return credentials, nil
				}
			} else if cerr != gooredis.ErrNil {
				goolog.ErrorF("sts redis get %s error: %s", key, cerr.Error())
			}
		}

		res, err := STSCredential(cfg)
		if err != nil {
			return nil, err
		}
		if res.Credentials == nil {
			return nil, fmt.Errorf("sts credentials is nil")
		}

		if redis != nil {
			b, merr := json.Marshal(res.Credentials)
			if merr == nil {
				ttl := time.Duration(cfg.GetExpire()) * time.Second
				if ttl > 10*time.Minute {
					ttl -= 5 * time.Minute
				} else if ttl > time.Minute {
					ttl -= time.Minute
				}
				if serr := redis.Set(key, string(b), ttl).Err(); serr != nil {
					goolog.ErrorF("sts redis set %s error: %s", key, serr.Error())
				}
			}
		}

		return res.Credentials, nil
	})

	if err != nil {
		return nil, err
	}

	cred, ok := result.(*sts.Credentials)
	if !ok || cred == nil {
		return nil, fmt.Errorf("sts credentials is nil")
	}
	return cred, nil
}
