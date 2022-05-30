package goo_oss

import "github.com/aliyun/aliyun-oss-go-sdk/oss"

var __oss *uploader

func Init(conf Config) {
	__oss, _ = New(conf)
}

func ContentType(value string) *uploader {
	return __oss.ContentType(value)
}

func Options(opts ...oss.Option) *uploader {
	return __oss.Options(opts...)
}

func Upload(filename string, body []byte) (string, error) {
	return __oss.Upload(filename, body)
}
