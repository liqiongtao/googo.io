package goo_file

import (
	"errors"
	"fmt"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"os"
	"sort"
)

var (
	maxLine = 1000000
)

func FileSort(filename, sortedFile string) (err error) {
	if !Exist(filename) {
		err = errors.New("文件不存在")
		return
	}

	var (
		// 部分文件
		partFiles []string

		// 第几部分
		partNum int

		// 临时数据
		data []string
	)

	defer func() {
		for _, file := range partFiles {
			if Exist(file) {
				os.Remove(file)
			}
		}
	}()

	defer func() {
		if err != nil {
			return
		}

		l := len(partFiles)
		if l == 0 {
			return
		}

		if l == 1 {
			os.Rename(partFiles[0], sortedFile)
			return
		}

		err = FileMerge(sortedFile, partFiles)
	}()

	err = ReadByLine(filename, func(b []byte, end bool) (err error) {
		defer func() {
			if l := len(data); l < maxLine && !end {
				return
			}

			var (
				partFile = fmt.Sprintf("%s.%d", filename, partNum)
				fh       *os.File
			)

			fh, err = os.OpenFile(partFile, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0755)
			if err != nil {
				goo_log.Error(err)
				return
			}
			defer fh.Close()

			sort.Strings(data)

			for _, s := range data {
				fh.WriteString(s)
			}

			goo_log.DebugF("产生一个文件：%s", partFile)

			partFiles = append(partFiles, partFile)

			data = []string{}
			partNum++
		}()

		data = append(data, string(b))

		return
	})

	return
}
