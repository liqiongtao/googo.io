package goo_oss

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"io"
	"os"
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

func (o *uploader) Upload(filename string, r io.Reader) (string, error) {
	var options []oss.Option

	if strings.Contains(filename, ".js") {
		options = append(options, oss.ContentType("application/javascript"))
	} else if strings.Contains(filename, ".css") {
		options = append(options, oss.ContentType("text/css"))
	} else if strings.Contains(filename, ".html") {
		options = append(options, oss.CacheControl("no-store"))
		options = append(options, oss.SetHeader("Pragma", "no-cache"))
	}

	if err := o.bucket.PutObject(filename, r, options...); err != nil {
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

func (o *uploader) UploadFile(filename, filepath string) (string, error) {
	if _, err := os.Stat(filepath); err != nil {
		goo_log.Error(err.Error())
		return "", err
	}

	f, err := os.Open(filepath)
	if err != nil {
		goo_log.Error(err.Error())
		return "", err
	}
	defer f.Close()

	return o.Upload(filename, f)
}

func (o *uploader) getClient() (*oss.Client, error) {
	return oss.New(o.conf.Endpoint, o.conf.AccessKeyId, o.conf.AccessKeySecret)
}

func (o *uploader) getBucket() (*oss.Bucket, error) {
	return o.client.Bucket(o.conf.Bucket)
}
