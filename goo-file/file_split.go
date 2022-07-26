package goo_file

import (
	"errors"
	"fmt"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
	"os"
	"runtime"
	"strings"
	"sync"
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

		index = strings.LastIndex(filename, ".")

		wg sync.WaitGroup
		ch = make(chan struct{}, runtime.NumCPU()*2)
	)

	err = ReadByLine(filename, func(b []byte, end bool) (err error) {
		defer func() {
			if l := len(data); l < maxLine && !end {
				return
			}

			wg.Add(1)
			ch <- struct{}{}

			func(partNum int, data []string) {
				goo_utils.AsyncFunc(func() {
					defer wg.Done()
					defer func() { <-ch }()

					var (
						partFile = fmt.Sprintf("%s.%d.%s", filename[:index], partNum, filename[index+1:])
						fh       *os.File
					)

					fh, err = os.OpenFile(partFile, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0755)
					if err != nil {
						goo_log.Error(err)
						return
					}
					defer fh.Close()

					for _, s := range data {
						fh.WriteString(s)
					}

					goo_log.DebugF("产生一个文件：%s", partFile)

					files = append(files, partFile)
				})
			}(partNum, data)

			data = []string{}
			partNum++
		}()

		data = append(data, string(b))
		return
	})

	wg.Wait()
	return
}
