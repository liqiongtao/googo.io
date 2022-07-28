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
func FileMerge(filename string, fileHandlers []*os.File) (err error) {
	defer func() {
		if err != nil {
			if Exist(filename + ".0") {
				os.Remove(filename + ".0")
			}
			return
		}

		os.Rename(filename+".0", filename)
	}()

	var fh *os.File

	fh, err = os.OpenFile(filename+".0", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0755)
	if err != nil {
		goo_log.Error(err)
		return
	}
	defer fh.Close()

	// 定义 分隔文件读句柄
	var rs = make([]*bufio.Reader, len(fileHandlers))
	for n, f := range fileHandlers {
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

	var strs []string

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
