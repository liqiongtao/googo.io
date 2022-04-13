package goo_upload

import (
	"bytes"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
	"log"
	"path"
)

var OSS *gooOSS

type OSSConfig struct {
	AccessKeyId     string `yaml:"access_key_id"`
	AccessKeySecret string `yaml:"access_key_secret"`
	Endpoint        string `yaml:"endpoint"`
	Bucket          string `yaml:"bucket"`
	Domain          string `yaml:"domain"`
}

func InitOSS(config OSSConfig) {
	var err error
	OSS, err = NewOSS(config)
	if err != nil {
		log.Panic(err.Error())
	}
}

func NewOSS(config OSSConfig) (*gooOSS, error) {
	o := &gooOSS{
		Config:  config,
		options: []oss.Option{},
	}

	client, err := o.getClient()
	if err != nil {
		goo_log.Error(err.Error())
		return nil, err
	}

	o.Client = client

	bucket, err := o.getBucket()
	if err != nil {
		goo_log.Error(err.Error())
		return nil, err
	}

	o.Bucket = bucket

	return o, nil
}

type gooOSS struct {
	Config  OSSConfig
	Client  *oss.Client
	Bucket  *oss.Bucket
	options []oss.Option
}

func (o *gooOSS) ContentType(value string) *gooOSS {
	o.options = append(o.options, oss.ContentType(value))
	return o
}

func (o *gooOSS) Options(opts ...oss.Option) *gooOSS {
	o.options = append(o.options, opts...)
	return o
}

func (o *gooOSS) Upload(filename string, body []byte) (string, error) {
	md5str := goo_utils.MD5(body)
	filename = fmt.Sprintf("%s/%s/%s%s", md5str[0:2], md5str[2:4], md5str[8:24], path.Ext(filename))

	if err := o.Bucket.PutObject(filename, bytes.NewReader(body), o.options...); err != nil {
		goo_log.Error(err.Error())
		return "", err
	}

	if o.Config.Domain != "" {
		return o.Config.Domain + filename, nil
	}

	url := "https://" + o.Config.Bucket + "." + o.Config.Endpoint + "/" + filename
	return url, nil
}

func (o *gooOSS) getClient() (*oss.Client, error) {
	return oss.New(o.Config.Endpoint, o.Config.AccessKeyId, o.Config.AccessKeySecret)
}

func (o *gooOSS) getBucket() (*oss.Bucket, error) {
	return o.Client.Bucket(o.Config.Bucket)
}
