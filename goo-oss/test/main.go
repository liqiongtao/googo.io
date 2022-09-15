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

var (
	AccessKeyId     = flag.String("access_key_id", "", "")
	AccessKeySecret = flag.String("access_key_secret", "", "")
	Endpoint        = flag.String("endpoint", "", "")
	Bucket          = flag.String("bucket", "", "")
	Domain          = flag.String("domain", "", "")
)

func main() {
	goo.FlagInit()

	args := os.Args
	if l := len(args); l < 2 {
		fmt.Println("请选择上传文件!")
		return
	}

	conf := goo_oss.Config{
		AccessKeyId:     *AccessKeyId,
		AccessKeySecret: *AccessKeySecret,
		Endpoint:        *Endpoint,
		Bucket:          *Bucket,
		Domain:          *Domain,
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

	for n, i := range args {
		if n == 0 {
			continue
		}

		b, err := ioutil.ReadFile(i)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		index := strings.LastIndex(i, "/")

		var filename string
		if index == -1 {
			filename = i
		} else {
			filename = i[index+1:]
		}

		md5 := goo_utils.MD5(b)
		filename = fmt.Sprintf("%s/%s/%s", md5[0:2], md5[2:4], filename)

		url, err := up.Upload(strings.ToLower(filename), b)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		fmt.Println(url)
	}
}
