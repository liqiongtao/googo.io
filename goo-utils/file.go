package goo_utils

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
)

const (
	EL = "\n"
)

// 文件名
func FILE() string {
	_, file, _, _ := runtime.Caller(1)
	return file
}

// 行号
func LINE() int {
	_, _, line, _ := runtime.Caller(1)
	return line
}

// 目录名称
func DIR() string {
	_, file, _, _ := runtime.Caller(1)
	return path.Dir(file) + "/"
}

// 写文件，支持路径创建
func WriteToFile(filename string, b []byte) error {
	dirname := path.Dir(filename)
	if _, err := os.Stat(dirname); err != nil {
		os.MkdirAll(dirname, 0755)
	}
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.Write(b); err != nil {
		return err
	}
	return nil
}

// 追踪信息
func Trace(skip int, trimPaths ...string) (arr []string) {
	arr = []string{}
	l := len(trimPaths)
	if skip == 0 {
		skip = 1
	}
	for i := skip; i < 16; i++ {
		_, file, line, _ := runtime.Caller(i)
		if file == "" {
			continue
		}
		if strings.Contains(file, ".pb.go") ||
			strings.Contains(file, "runtime/") ||
			strings.Contains(file, "src/testing") ||
			strings.Contains(file, "pkg/mod/") ||
			strings.Contains(file, "vendor/") {
			continue
		}
		if l > 0 {
			for _, trimPath := range trimPaths {
				file = strings.Replace(file, trimPath, "", -1)
			}
		}
		arr = append(arr, fmt.Sprintf("%s %dL", file, line))
	}
	return
}
