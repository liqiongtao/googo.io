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

func Trace(skip int) []string {
	trace := []string{}
	if skip == 0 {
		skip = 2
	}
	for i := skip; i < 16; i++ {
		_, file, line, _ := runtime.Caller(i)
		fmt.Println("--1--", file, line)
		if file == "" ||
			strings.Index(file, "runtime") > 0 ||
			strings.Index(file, "src/testing") > 0 ||
			strings.Index(file, "pkg/mod") > 0 ||
			strings.Index(file, "vendor") > 0 {
			continue
		}
		trace = append(trace, fmt.Sprintf("%s %dL", file, line))
	}
	return trace
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
