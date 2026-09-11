package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/liqiongtao/googo.io/goo"
	goofile "github.com/liqiongtao/googo.io/goo-file"
	goooss "github.com/liqiongtao/googo.io/goo-oss"
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

	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("请选择上传文件!")
		return
	}

	conf := goooss.Config{
		AccessKeyId:     *AccessKeyIdFlag,
		AccessKeySecret: *AccessKeySecretFlag,
		Endpoint:        *EndpointFlag,
		Bucket:          *BucketFlag,
		Domain:          *DomainFlag,
	}

	up, err := goooss.New(conf)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, filePath := range args {
		_, err := os.Stat(filePath)
		if err != nil {
			fmt.Println(err)
			continue
		}

		md5str, err := goofile.MD5(filePath)
		if err != nil {
			fmt.Println(err)
			continue
		}

		var filename string
		{
			index := strings.LastIndex(filePath, "/")
			nw := time.Now()
			if index == -1 {
				filename = fmt.Sprintf("%s/%s/%s/%s", nw.Format("2006"), md5str[0:2], md5str[2:4], filePath)
			} else {
				filename = fmt.Sprintf("%s/%s/%s/%s", nw.Format("2006"), md5str[0:2], md5str[2:4], filePath[index+1:])
			}
		}

		{
			index := strings.LastIndex(filename, ".")
			if index > 0 {
				filename = fmt.Sprintf("%s_%s.%s", filename[:index], md5str[8:24], filename[index+1:])
			} else {
				filename = fmt.Sprintf("%s_%s", filename, md5str[8:24])
			}
		}

		url, err := up.UploadFile(filename, filePath)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println(url)
	}
}
