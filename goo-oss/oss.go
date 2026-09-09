package goo_oss

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"sync"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
)

var (
	__oss *Uploader
	__mu  sync.RWMutex
)

func Init(conf Config) (err error) {
	o, err := New(conf)
	if err != nil {
		return err
	}
	__mu.Lock()
	__oss = o
	__mu.Unlock()
	return
}

func requireOSS() (*Uploader, error) {
	__mu.RLock()
	o := __oss
	__mu.RUnlock()
	if o == nil {
		return nil, errors.New("oss not initialized, call goo_oss.Init first")
	}
	return o, nil
}

func Default() *Uploader {
	__mu.RLock()
	defer __mu.RUnlock()
	return __oss
}

func Client() (*oss.Client, error) {
	o, err := requireOSS()
	if err != nil {
		return nil, err
	}
	if o.Client == nil {
		return nil, errors.New("oss client is nil")
	}
	return o.Client, nil
}

func Bucket() (*oss.Bucket, error) {
	o, err := requireOSS()
	if err != nil {
		return nil, err
	}
	if o.Bucket != nil {
		return o.Bucket, nil
	}
	bucket, err := o.Client.Bucket(o.conf.Bucket)
	if err != nil {
		return nil, err
	}
	return bucket, nil
}

func ContentType(value string) (*Uploader, error) {
	o, err := requireOSS()
	if err != nil {
		return nil, err
	}
	// 返回独立副本，避免污染全局单例 options
	return &Uploader{
		conf:    o.conf,
		Client:  o.Client,
		Bucket:  o.Bucket,
		options: []oss.Option{oss.ContentType(value)},
	}, nil
}

func Options(opts ...oss.Option) (*Uploader, error) {
	o, err := requireOSS()
	if err != nil {
		return nil, err
	}
	copied := make([]oss.Option, len(opts))
	copy(copied, opts)
	return &Uploader{
		conf:    o.conf,
		Client:  o.Client,
		Bucket:  o.Bucket,
		options: copied,
	}, nil
}

func Upload(filename string, r io.Reader) (string, error) {
	o, err := requireOSS()
	if err != nil {
		return "", err
	}
	return o.Upload(filename, r)
}

func UploadFile(filename, filepath string) (string, error) {
	o, err := requireOSS()
	if err != nil {
		return "", err
	}
	return o.UploadFile(filename, filepath)
}

func GetAppendPosition(objectKey string) (int64, error) {
	o, err := requireOSS()
	if err != nil {
		return 0, err
	}
	hd, err := o.Bucket.GetObjectDetailedMeta(objectKey)
	if err != nil {
		var v oss.ServiceError
		if errors.As(err, &v) {
			switch v.Code {
			case "NoSuchKey":
				return 0, nil
			}
		}

		goo_log.Error(err.Error())
		return 0, err
	}

	if hd == nil {
		goo_log.Error("httpHeader is nil")
		return 0, fmt.Errorf("httpHeader is nil")
	}

	position, err := strconv.ParseInt(hd.Get("x-oss-next-append-position"), 10, 64)
	if err != nil {
		goo_log.Error(err.Error())
		return 0, err
	}

	return position, nil
}

// AppendObject 追加写入；遇 PositionNotEqualToLength 时重新取 position 重试，缓解并发追加冲突。
func AppendObject(objectKey string, b []byte) (int64, error) {
	o, err := requireOSS()
	if err != nil {
		return 0, err
	}

	const maxRetry = 5
	var lastErr error
	for i := 0; i < maxRetry; i++ {
		appendPosition, err := GetAppendPosition(objectKey)
		if err != nil {
			return 0, err
		}
		next, err := o.Bucket.AppendObject(objectKey, bytes.NewReader(b), appendPosition)
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
