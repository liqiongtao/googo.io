package goooss

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"path"
	"strconv"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	goolog "github.com/liqiongtao/googo.io/goo-log"
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
	return __oss.Bucket
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

func applyPrefix(prefix, objectKey string) string {
	objectKey = strings.TrimPrefix(objectKey, "/")
	prefix = strings.Trim(prefix, "/")
	if prefix == "" {
		return objectKey
	}
	if objectKey == prefix || strings.HasPrefix(objectKey, prefix+"/") {
		return objectKey
	}
	return path.Join(prefix, objectKey)
}

func GetAppendPosition(objectKey string) (int64, error) {
	objectKey = applyPrefix(__oss.conf.Prefix, objectKey)
	hd, err := __oss.Bucket.GetObjectDetailedMeta(objectKey)
	if err != nil {
		var v oss.ServiceError
		if errors.As(err, &v) {
			switch v.Code {
			case "NoSuchKey":
				return 0, nil
			}
		}

		goolog.Error(err.Error())
		return 0, err
	}

	if hd == nil {
		goolog.Error("httpHeader is nil")
		return 0, fmt.Errorf("httpHeader is nil")
	}

	position, err := strconv.ParseInt(hd.Get("x-oss-next-append-position"), 10, 64)
	if err != nil {
		goolog.Error(err.Error())
		return 0, err
	}

	return position, nil
}

// AppendObject 追加写入；遇 PositionNotEqualToLength 时重新取 position 重试，缓解并发追加冲突。
func AppendObject(objectKey string, b []byte) (int64, error) {
	const maxRetry = 5
	var lastErr error
	objectKey = applyPrefix(__oss.conf.Prefix, objectKey)
	for i := 0; i < maxRetry; i++ {
		appendPosition, err := GetAppendPosition(objectKey)
		if err != nil {
			return 0, err
		}
		next, err := __oss.Bucket.AppendObject(objectKey, bytes.NewReader(b), appendPosition)
		if err == nil {
			return next, nil
		}
		lastErr = err
		var se oss.ServiceError
		if errors.As(err, &se) && se.Code == "PositionNotEqualToLength" {
			continue
		}
		return 0, err
	}
	return 0, lastErr
}
