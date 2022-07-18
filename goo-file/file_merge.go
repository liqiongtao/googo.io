package goo_file

import (
	"bufio"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"io"
	"os"
	"sort"
	"strings"
)

// 把多个文件内容合并到一个文件里面，合并时做好排序
func FileMerge(filename string, files []*os.File) (err error) {
	defer func() {
		if err != nil {
			if Exist(filename + ".0") {
				if err = os.Remove(filename + ".0"); err != nil {
					goo_log.Error(err)
				}
			}
			return
		}

		if err = os.Rename(filename+".0", filename); err != nil {
			goo_log.Error(err)
		}
	}()

	var fh *os.File

	fh, err = os.OpenFile(filename+".0", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0755)
	if err != nil {
		goo_log.Error(err)
		return
	}
	defer fh.Close()

	// 定义 分隔文件读句柄
	var rs = make([]*bufio.Reader, len(files))
	for n, f := range files {
		f.Seek(0, 0)
		rs[n] = bufio.NewReader(f)
	}

	// 定义 每个文件 拿到的一行字符串
	var data = map[string]int{}
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
		fh.WriteString(str)

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

	return
}
