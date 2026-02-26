package goo_oss

import (
	"bytes"
	"fmt"
	"io"
	"strconv"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
)

var __oss *Uploader

func Init(conf Config) (err error) {
	__oss, err = New(conf)
	return
}

func Default() *Uploader {
	return __oss
}

func Client() *oss.Client {
	return __oss.Client
}

func Bucket() *oss.Bucket {
	bucket, err := __oss.Client.Bucket(__oss.conf.Bucket)
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

func GetAppendPosition(objectKey string) (int64, error) {
	hd, err := __oss.Bucket.GetObjectDetailedMeta(objectKey)
	if err != nil {
		switch err.(oss.ServiceError).Code {
		case "NoSuchKey":
			return 0, nil
		default:
			return 0, err
		}
	}

	if hd == nil {
		return 0, fmt.Errorf("httpHeader is nil")
	}

	position, err := strconv.ParseInt(hd.Get("x-oss-next-append-position"), 10, 64)
	if err != nil {
		return 0, err
	}

	return position, nil
}

func AppendObject(objectKey string, b []byte) (int64, error) {
	appendPosition, err := GetAppendPosition(objectKey)
	if err != nil {
		return 0, err
	}
	return __oss.Bucket.AppendObject(objectKey, bytes.NewReader(b), appendPosition)
}
