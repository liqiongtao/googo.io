package goo_cos

import (
	"fmt"
	"net/url"
)

type StsConfig struct {
	SecretId  string `json:"secret_id" yaml:"secret_id"`   // 密钥ID
	SecretKey string `json:"secret_key" yaml:"secret_key"` // 密钥
	Appid     string `json:"appid" yaml:"appid"`           // 应用ID
	Bucket    string `json:"bucket" yaml:"bucket"`         // 存储桶名称
	Region    string `json:"region" yaml:"region"`         // 存储桶所在的地域

	Action []string `json:"action" yaml:"action"` // 执行权限
	Expire int64    `json:"expire" yaml:"expire"` // 有效时间，单位秒
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
	SecretId     string `json:"secret_id" yaml:"secret_id"`         // 密钥ID
	SecretKey    string `json:"secret_key" yaml:"secret_key"`       // 密钥
	SessionToken string `json:"session_token" yaml:"session_token"` // 临时key时，需要传入
	Bucket       string `json:"bucket" yaml:"bucket"`               // 存储桶名称
	Region       string `json:"region" yaml:"region"`               // 存储桶所在的地域
	BaseDir      string `json:"base_dir" yaml:"base_dir"`           // 基础目录
}

func (c CosConfig) GetBucketURL() (*url.URL, error) {
	str := fmt.Sprintf("https://%s.cos.%s.myqcloud.com", c.Bucket, c.Region)
	return url.Parse(str)
}
