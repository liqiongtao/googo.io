package goo_oss

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"io"
	"os"
	"strings"
	"time"
)

type Uploader struct {
	conf    Config
	Client  *oss.Client
	Bucket  *oss.Bucket
	options []oss.Option
}

func New(conf Config) (*Uploader, error) {
	o := &Uploader{
		conf:    conf,
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

func (o *Uploader) ContentType(value string) *Uploader {
	o.options = append(o.options, oss.ContentType(value))
	return o
}

func (o *Uploader) Options(opts ...oss.Option) *Uploader {
	o.options = append(o.options, opts...)
	return o
}

func (o *Uploader) Upload(filename string, r io.Reader) (string, error) {
	var options []oss.Option

	if strings.Contains(filename, ".js") {
		options = append(options, oss.ContentType("application/javascript"))
	} else if strings.Contains(filename, ".css") {
		options = append(options, oss.ContentType("text/css"))
	} else if strings.Contains(filename, ".html") {
		options = append(options, oss.CacheControl("no-store"))
		options = append(options, oss.SetHeader("Pragma", "no-cache"))
	}

	// 拼接前缀
	filename = fmt.Sprintf("%s/%s", o.conf.Prefix, filename)
	filename = strings.ReplaceAll(filename, "///", "/")
	filename = strings.ReplaceAll(filename, "//", "/")

	for i := 0; i < 3; i++ {
		err := o.Bucket.PutObject(filename, r, options...)
		if err == nil {
			break
		}

		goo_log.Error(err.Error())

		if i+1 == 3 {
			return "", err
		}

		time.Sleep(time.Second)
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

func (o *Uploader) UploadFile(filename, filepath string) (string, error) {
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

func (o *Uploader) getClient() (*oss.Client, error) {
	return oss.New(o.conf.Endpoint, o.conf.AccessKeyId, o.conf.AccessKeySecret)
}

func (o *Uploader) getBucket() (*oss.Bucket, error) {
	return o.Client.Bucket(o.conf.Bucket)
}
