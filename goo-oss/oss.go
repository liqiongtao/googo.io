package goo_oss

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"io"
)

var __oss *Uploader

func Init(conf Config) {
	__oss, _ = New(conf)
}

func Client() *oss.Client {
	return __oss.client
}

func Bucket() *oss.Bucket {
	bucket, err := __oss.client.Bucket(__oss.conf.Bucket)
	if err != nil {
		goo_log.Error(err)
	}
	return bucket
}

func ContentType(value string) *Uploader {
	return __oss.ContentType(value)
}

func Options(opts ...oss.Option) *Uploader {
	return __oss.Options(opts...)
}

func Upload(filename string, r io.Reader) (string, error) {
	return __oss.Upload(filename, r)
}

func UploadFile(filename, filepath string) (string, error) {
	return __oss.UploadFile(filename, filepath)
}
