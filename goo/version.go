package goo

import (
	"flag"
)

/**
  - flag.go

	func init() {
		if *goo.VersionFlag {
			fmt.Println(goo.Version)
			os.Exit(0)
		}
	}

  - deploy.sh

	// 定义
	version=v1.0.1
	buildVersion="${version}.$(date +%Y%m%d).$(date +%H%M)"

	// 编译
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -X goo.Version=$buildVersion" -o ss

  - 执行

	./ss -v
*/

var (
	Version     string
	VersionFlag = flag.Bool("v", false, "version")
)
