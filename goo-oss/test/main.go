package main

import (
	"flag"
	"fmt"
	"github.com/liqiongtao/googo.io/goo"
	goo_file "github.com/liqiongtao/googo.io/goo-file"
	goo_oss "github.com/liqiongtao/googo.io/goo-oss"
	"os"
	"strings"
	"time"
)

// go build -ldflags "-s -w" -o oss

var (
	AccessKeyIdFlag     = flag.String("access_key_id", "", "")
	AccessKeySecretFlag = flag.String("access_key_secret", "", "")
	EndpointFlag        = flag.String("endpoint", "", "")
	BucketFlag          = flag.String("bucket", "", "")
	DomainFlag          = flag.String("domain", "", "")
)

func main() {
	goo.FlagInit()

	args := os.Args
	if l := len(args); l < 2 {
		fmt.Println("请选择上传文件!")
		return
	}

	conf := goo_oss.Config{
		AccessKeyId:     *AccessKeyIdFlag,
		AccessKeySecret: *AccessKeySecretFlag,
		Endpoint:        *EndpointFlag,
		Bucket:          *BucketFlag,
		Domain:          *DomainFlag,
	}

	if conf.AccessKeyId == "" {
		conf.AccessKeyId = ""
	}
	if conf.AccessKeySecret == "" {
		conf.AccessKeySecret = ""
	}
	if conf.Endpoint == "" {
		conf.Endpoint = ""
	}
	if conf.Bucket == "" {
		conf.Bucket = ""
	}
	if conf.Domain == "" {
		conf.Domain = ""
	}

	up, err := goo_oss.New(conf)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for i := 1; i < len(args); i++ {
		_, err := os.Stat(args[i])
		if err != nil {
			fmt.Println(err)
			continue
		}

		md5str, err := goo_file.MD5(args[i])
		if err != nil {
			fmt.Println(err)
			continue
		}

		var filename string
		{
			index := strings.LastIndex(args[i], "/")
			nw := time.Now()
			if index == -1 {
				filename = fmt.Sprintf("%s/%s/%s/%s", nw.Format("2006"), md5str[0:2], md5str[2:4], args[i])
			} else {
				filename = fmt.Sprintf("%s/%s/%s/%s", nw.Format("2006"), md5str[0:2], md5str[2:4], args[i][index+1:])
			}
		}

		{
			index := strings.Index(filename, ".")
			filename = fmt.Sprintf("%s_%s.%s", filename[:index], md5str[8:24], filename[index+1:])
		}

		url, err := up.UploadFile(filename, args[i])
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println(url)
	}
}
