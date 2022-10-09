package main

import (
	"flag"
	"fmt"
	"github.com/liqiongtao/googo.io/goo"
	goo_oss "github.com/liqiongtao/googo.io/goo-oss"
	"io/ioutil"
	"os"
)

// go build -ldflags "-s -w" -o oss

var (
	AccessKeyIdFlag     = flag.String("access_key_id", "", "")
	AccessKeySecretFlag = flag.String("access_key_secret", "", "")
	EndpointFlag        = flag.String("endpoint", "oss-cn-beijing.aliyuncs.com", "")
	BucketFlag          = flag.String("bucket", "video-ai", "")
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

	url, err := up.Upload(args[1], b)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(url)
}
