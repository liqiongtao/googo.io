package goo_cos

import (
	"fmt"
	"net/url"
)

type StsConfig struct {
	SecretId  string // 密钥ID
	SecretKey string // 密钥
	Appid     string // 应用ID
	Bucket    string // 存储桶名称
	Region    string // 存储桶所在的地域

	Action []string // 执行权限
	Expire int64    // 有效时间，单位秒
}

func (c StsConfig) GetAction() []string {
	if c.Action == nil {
		return DefaultAction
	}
	return c.Action
}

func (c StsConfig) GetExpire() int64 {
	if c.Expire == 0 {
		return DurationSeconds
	}
	return c.Expire
}

func (c StsConfig) GetResource() []string {
	return []string{
		fmt.Sprintf(ResourceTemplate, c.Region, c.Appid, c.Bucket),
	}
}

type CosConfig struct {
	SecretId     string // 密钥ID
	SecretKey    string // 密钥
	SessionToken string // 临时key时，需要传入
	Bucket       string // 存储桶名称
	Region       string // 存储桶所在的地域
	Dir          string // 存储目录
}

func (c CosConfig) GetBucketURL() (*url.URL, error) {
	str := fmt.Sprintf("https://%s.cos.%s.myqcloud.com", c.Bucket, c.Region)
	return url.Parse(str)
}
