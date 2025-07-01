package goo_cos

import (
	"fmt"
	"net/url"
)

type StsConfig struct {
	SecretId  string
	SecretKey string
	Appid     string
	Bucket    string
	Region    string

	// 临时密钥动作
	Action []string
	// 临时密钥有效时间，单位秒
	Expire int64
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
	SecretId     string
	SecretKey    string
	SessionToken string // 临时key时，需要传入
	Bucket       string
	Region       string
}

func (c CosConfig) GetBucketURL() (*url.URL, error) {
	str := fmt.Sprintf("https://%s.cos.%s.myqcloud.com", c.Bucket, c.Region)
	return url.Parse(str)
}
