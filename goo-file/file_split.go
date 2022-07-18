package goo_file

import (
	"errors"
	"fmt"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"os"
	"strings"
)

func FileSplit(filename string, maxLine int) (files []string, err error) {
	files = []string{}

	if !Exist(filename) {
		err = errors.New("文件不存在")
		return
	}

	var (
		partNum int
		data    []string
		index   = strings.LastIndex(filename, ".")
	)

	err = ReadByLine(filename, func(b []byte, end bool) (err error) {
		defer func() {
			if l := len(data); l < maxLine && !end {
				return
			}

			var (
				partFile = fmt.Sprintf("%s.%d.%s", filename[:index], partNum, filename[index+1:])
				fh       *os.File
			)

			fh, err = os.OpenFile(partFile, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0755)
			if err != nil {
				goo_log.Error(err)
				return
			}

			for _, s := range data {
				fh.WriteString(s)
			}

			fh.Close()

			files = append(files, partFile)

			data = []string{}
			partNum++
		}()

		data = append(data, string(b))

		return
	})

	return
}
