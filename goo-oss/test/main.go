package main

import (
	"flag"
	"fmt"
	"github.com/liqiongtao/googo.io/goo"
	goo_oss "github.com/liqiongtao/googo.io/goo-oss"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
	"io/ioutil"
	"os"
	"strings"
)

// go build -ldflags "-s -w" -o oss

var (
	AccessKeyIdFlag     = flag.String("access_key_id", "", "")
	AccessKeySecretFlag = flag.String("access_key_secret", "", "")
	EndpointFlag        = flag.String("endpoint", "", "")
	BucketFlag          = flag.String("bucket", "", "")
	DomainFlag          = flag.String("domain", "", "")
	filenameFlag        = flag.String("filename", "", "")
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
		conf.Endpoint = "oss-cn-beijing.aliyuncs.com"
	}
	if conf.Bucket == "" {
		conf.Bucket = ""
	}

	up, err := goo_oss.New(conf)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	b, err := ioutil.ReadFile(args[1])
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var filename string
	{
		index := strings.LastIndex(args[1], "/")
		if index == -1 {
			filename = args[1]
		} else {
			filename = args[1][index+1:]
		}
	}

	md5 := goo_utils.MD5(b)

	{
		index := strings.LastIndex(filename, ".")
		filename = fmt.Sprintf("%s_%s.%s", filename[:index], md5[8:24], filename[index+1:])
	}

	filename = fmt.Sprintf("%s/%s/%s", md5[0:2], md5[2:4], filename)

	if *filenameFlag != "" {
		filename = *filenameFlag
	}

	url, err := up.Upload(strings.ToLower(filename), b)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(url)
}
