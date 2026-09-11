package goooss

import (
	"fmt"
	"strings"
)

type Config struct {
	AccessKeyId     string `json:"access_key_id" yaml:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret" yaml:"access_key_secret"`
	Endpoint        string `json:"endpoint" yaml:"endpoint"`
	Bucket          string `json:"bucket" yaml:"bucket"`
	Domain          string `json:"domain" yaml:"domain"`
	Prefix          string `json:"prefix" yaml:"prefix"`
}

func (c Config) BaseUrl() string {
	endpoint := strings.TrimPrefix(c.Endpoint, "https://")
	endpoint = strings.TrimPrefix(endpoint, "http://")
	return fmt.Sprintf("https://%s.%s", c.Bucket, endpoint)
}
