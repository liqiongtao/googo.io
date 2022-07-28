package goo_file

import (
	"bufio"
	"fmt"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
	"io"
	"os"
	"runtime"
	"sort"
	"strings"
	"sync"
)

var (
	tempMergeFiles []string
)

func FileMerge(file string, files []string) (err error) {
	defer func() {
		for _, _file := range tempMergeFiles {
			if Exist(_file) {
				os.Remove(_file)
			}
		}
	}()

	var (
		index int
		size  = 5
	)

	for {
		l := len(files)
		if l == 0 {
			return
		}
		if l == 1 {
			if _file := files[0]; Exist(_file) {
				os.Rename(_file, file)
			}
			return
		}

		filename := fmt.Sprintf("%s.%d", file, index)
		filesArr := goo_utils.SplitStringArray(files, size)

		files = fileGroupMerge(filename, filesArr)

		index++
	}
}

func fileGroupMerge(file string, filesArr [][]string) (files []string) {
	files = []string{}

	var (
		wg sync.WaitGroup
		ch = make(chan struct{}, runtime.NumCPU()/2)
	)

	for n, _files := range filesArr {
		l := len(_files)
		if l == 0 {
			continue
		}
		if l == 1 {
			files = append(files, _files[0])
			continue
		}

		wg.Add(1)
		ch <- struct{}{}

		_file := fmt.Sprintf("%s.%d", file, n)

		files = append(files, _file)
		tempMergeFiles = append(tempMergeFiles, _file)

		goo_log.DebugF("文件合并，临时文件: %s", _file)

		func(_file string, _files []string) {
			goo_utils.AsyncFunc(func() {
				defer wg.Done()
				defer func() { <-ch }()

				fileMergeHandler(_file, _files)
			})
		}(_file, _files)
	}

	wg.Wait()

	return
}

func fileMergeHandler(file string, files []string) (err error) {
	var fh *os.File

	fh, err = os.OpenFile(file, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0755)
	if err != nil {
		return
	}

	var (
		handlers []*os.File
		rs       []*bufio.Reader
		wg       sync.WaitGroup
	)

	for _, _file := range files {
		wg.Add(1)

		func(_file string) {
			goo_utils.AsyncFunc(func() {
				defer wg.Done()

				var f *os.File

				f, err = os.OpenFile(_file, os.O_RDWR, 0755)
				if err != nil {
					goo_log.Error(err)
					return
				}

				rs = append(rs, bufio.NewReader(f))
				handlers = append(handlers, f)
			})
		}(_file)
	}

	wg.Wait()

	defer func() {
		for _, f := range handlers {
			if f != nil {
				f.Close()
			}
		}
	}()

	var (
		data = map[string]int{}
		strs []string
	)

	for n, r := range rs {
		s, _ := r.ReadString('\n')
		if strings.TrimSpace(s) == "" {
			continue
		}
		data[s] = n
	}

	for {
		if l := len(data); l == 0 {
			break
		}

		// 收集每个文件的行字符串
		var arr []string
		for s := range data {
			arr = append(arr, s)
		}

		// 排序
		sort.Strings(arr)

		// 取第一个字符串
		str := arr[0]

		// 写入到文件
		{
			strs = append(strs, str)
			if l := len(strs); l >= 1000 {
				fh.WriteString(strings.Join(strs, ""))
				strs = []string{}
			}
		}

		// 该字符串对应的索引
		n := data[str]

		// 删除已经使用的str
		delete(data, str)

		var s string
		s, err = rs[n].ReadString('\n')
		if err != nil {
			if io.EOF == err {
				err = nil
				continue
			}
			goo_log.Error(err)
			return
		}

		if strings.TrimSpace(s) == "" {
			continue
		}

		data[s] = n
	}

	if l := len(strs); l > 0 {
		fh.WriteString(strings.Join(strs, ""))
	}

	return
}
