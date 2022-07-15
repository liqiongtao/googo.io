package goo_file

import (
	"errors"
	"fmt"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"os"
	"sort"
	"strings"
)

var (
	maxLine = 1000000
)

func FileSort(filename string) (err error) {
	if !Exist(filename) {
		err = errors.New("文件不存在")
		return
	}

	index := strings.LastIndex(filename, ".")
	sortedFile := fmt.Sprintf("%s.sort.%s", filename[:index], filename[index+1:])

	var (
		// 部分文件
		partFiles []string

		// 部分文件句柄
		partFileHandlers []*os.File

		// 第几部分
		partNum int

		// 临时数据
		data []string
	)

	defer func() {
		for _, file := range partFiles {
			if !Exist(file) {
				return
			}
			os.Remove(file)
		}
	}()

	defer func() {
		l := len(partFiles)
		if l == 0 {
			return
		}
		if l == 1 {
			os.Rename(partFiles[0], sortedFile)
			return
		}
		err = FileMerge(sortedFile, partFileHandlers)
	}()

	return ReadByLine(filename, func(b []byte, end bool) (err error) {
		defer func() {
			if l := len(data); l < maxLine && !end {
				return
			}

			var (
				partFile = fmt.Sprintf("%s.sort.%d.%s", filename[:index], partNum, filename[index+1:])
				fh       *os.File
			)

			fh, err = os.OpenFile(partFile, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0755)
			if err != nil {
				goo_log.Error(err)
				return
			}

			sort.Strings(data)

			for _, s := range data {
				fh.WriteString(s)
			}

			partFiles = append(partFiles, partFile)
			partFileHandlers = append(partFileHandlers, fh)

			data = []string{}
			partNum++
		}()

		data = append(data, string(b))

		return
	})
}
