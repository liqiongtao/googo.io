package goo_oss

import (
	"bytes"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
	"path"
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
	md5str := goo_utils.MD5(body)
	filepath := fmt.Sprintf("%s/%s/", md5str[0:2], md5str[2:4])

	if index := strings.LastIndexByte(filename, '/'); index > 0 {
		filename = filename[index+1:]
	}

	ext := path.Ext(filename)
	index := strings.Index(filename, ext)
	filename = filename[:index] + "_" + md5str[8:16] + filename[index:]

	if err := o.bucket.PutObject(filepath+filename, bytes.NewReader(body), o.options...); err != nil {
		goo_log.Error(err.Error())
		return "", err
	}

	if o.conf.Domain != "" {
		return o.conf.Domain + filepath + filename, nil
	}

	url := "https://" + o.conf.Bucket + "." + o.conf.Endpoint + "/" + filepath + filename
	return url, nil
}

func (o *uploader) getClient() (*oss.Client, error) {
	return oss.New(o.conf.Endpoint, o.conf.AccessKeyId, o.conf.AccessKeySecret)
}

func (o *uploader) getBucket() (*oss.Bucket, error) {
	return o.client.Bucket(o.conf.Bucket)
}
