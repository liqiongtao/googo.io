package goo_oss

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
	"strings"
)

type uploader struct {
	conf    Config
	client  *oss.Client
	bucket  *oss.Bucket
	options []oss.Option
}

func New(conf Config) (*uploader, error) {
	o := &uploader{
		conf:    conf,
		options: []oss.Option{},
	}

	client, err := o.getClient()
	if err != nil {
		goo_log.Error(err.Error())
		return nil, err
	}

	o.client = client

	bucket, err := o.getBucket()
	if err != nil {
		goo_log.Error(err.Error())
		return nil, err
	}

	o.bucket = bucket

	return o, nil
}

func (o *uploader) ContentType(value string) *uploader {
	o.options = append(o.options, oss.ContentType(value))
	return o
}

func (o *uploader) Options(opts ...oss.Option) *uploader {
	o.options = append(o.options, opts...)
	return o
}

func (o *uploader) Upload(filename string, body []byte) (string, error) {
	if filename == "" {
		return "", errors.New("文件名为空")
	}

	md5str := goo_utils.MD5(body)

	{
		index := strings.LastIndex(filename, "/")
		if index == -1 {
			filename = fmt.Sprintf("%s/%s/%s", md5str[0:2], md5str[2:4], filename)
		} else {
			filename = fmt.Sprintf("%s/%s/%s", md5str[0:2], md5str[2:4], filename[index+1:])
		}
	}

	{
		index := strings.Index(filename, ".")
		filename = fmt.Sprintf("%s_%s.%s", filename[:index], md5str[8:24], filename[index+1:])
	}

	filename = strings.ToLower(filename)

	if err := o.bucket.PutObject(filename, bytes.NewReader(body), o.options...); err != nil {
		goo_log.Error(err.Error())
		return "", err
	}

	if filename[0:1] != "/" {
		filename = "/" + filename
	}

	if o.conf.Domain != "" {
		if idx, l := strings.LastIndex(o.conf.Domain, "/"), len(o.conf.Domain); idx+1 == l {
			o.conf.Domain = o.conf.Domain[:l-1]
		}
		return o.conf.Domain + filename, nil
	}

	url := "https://" + o.conf.Bucket + "." + o.conf.Endpoint + filename
	return url, nil
}

func (o *uploader) getClient() (*oss.Client, error) {
	return oss.New(o.conf.Endpoint, o.conf.AccessKeyId, o.conf.AccessKeySecret)
}

func (o *uploader) getBucket() (*oss.Bucket, error) {
	return o.client.Bucket(o.conf.Bucket)
}
