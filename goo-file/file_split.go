package goo_file

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"

	goo_log "github.com/liqiongtao/googo.io/goo-log"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
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

		base, ext = splitFilename(filename)

		wg sync.WaitGroup
		mu sync.Mutex
		ch = make(chan struct{}, runtime.NumCPU()*2)
	)

	err = ReadByLine(filename, func(b []byte, end bool) (err error) {
		defer func() {
			if l := len(data); l < maxLine && !end {
				return
			}
			if len(data) == 0 {
				return
			}

			wg.Add(1)
			ch <- struct{}{}

			func(partNum int, data []string) {
				goo_utils.AsyncFunc(func() {
					defer wg.Done()
					defer func() { <-ch }()

					partFile := partFilename(base, ext, partNum)
					fh, openErr := os.OpenFile(partFile, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
					if openErr != nil {
						goo_log.Error(openErr)
						mu.Lock()
						if err == nil {
							err = openErr
						}
						mu.Unlock()
						return
					}
					defer fh.Close()

					for _, s := range data {
						if _, werr := fh.WriteString(s); werr != nil {
							goo_log.Error(werr)
							mu.Lock()
							if err == nil {
								err = werr
							}
							mu.Unlock()
							return
						}
					}

					goo_log.DebugF("产生一个文件: %s", partFile)

					mu.Lock()
					files = append(files, partFile)
					mu.Unlock()
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

func splitFilename(filename string) (base, ext string) {
	index := strings.LastIndex(filename, ".")
	if index < 0 {
		return filename, ""
	}
	return filename[:index], filename[index+1:]
}

func partFilename(base, ext string, partNum int) string {
	if ext == "" {
		return fmt.Sprintf("%s.%d", base, partNum)
	}
	return fmt.Sprintf("%s.%d.%s", base, partNum, ext)
}
