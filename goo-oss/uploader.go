package goo_oss

import (
	"io"
	"os"
	"path"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
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
	if o.conf.Prefix != "" {
		filename = path.Join(o.conf.Prefix, filename)
	}

	if err := o.Bucket.PutObject(filename, r, options...); err != nil {
		goo_log.Error("Oss Upload Failed", err.Error(), filename)
		return "", err
	}

	if o.conf.Domain != "" {
		domain := strings.TrimSuffix(o.conf.Domain, "/")
		return domain + "/" + filename, nil
	}

	url := "https://" + o.conf.Bucket + "." + o.conf.Endpoint + path.Join("/", filename)
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
	defer func() { _ = f.Close() }()

	return o.Upload(filename, f)
}

func (o *Uploader) getClient() (*oss.Client, error) {
	if strings.HasPrefix(o.conf.Endpoint, "http") {
		return oss.New(o.conf.Endpoint, o.conf.AccessKeyId, o.conf.AccessKeySecret)
	}
	return oss.New("https://"+o.conf.Endpoint, o.conf.AccessKeyId, o.conf.AccessKeySecret)
}

func (o *Uploader) getBucket() (*oss.Bucket, error) {
	return o.Client.Bucket(o.conf.Bucket)
}
