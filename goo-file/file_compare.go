package goo_file

import (
	"bufio"
	"errors"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"io"
	"os"
)

// 文件内容对比，输出增加的、减少的内容
// 文件内容要先做好排序
func Compare(srcFile, targetFile, appendFile, reduceFile string) (err error) {
	if !Exist(srcFile) || !Exist(targetFile) {
		err = errors.New("文件不存在")
		return
	}

	defer func() {
		if err != nil {
			if Exist(appendFile + ".0") {
				os.Remove(appendFile + ".0")
			}
			if Exist(reduceFile + ".0") {
				os.Remove(reduceFile + ".0")
			}
			return
		}

		os.Rename(appendFile+".0", appendFile)
		os.Rename(reduceFile+".0", reduceFile)
	}()

	var (
		f1, f2, f3, f4 *os.File
		r1, r2         *bufio.Reader
	)

	f1, err = os.OpenFile(srcFile, os.O_RDWR, 0755)
	if err != nil {
		goo_log.Error(err)
		return
	}
	defer f1.Close()

	f2, err = os.OpenFile(targetFile, os.O_RDWR, 0755)
	if err != nil {
		goo_log.Error(err)
		return
	}
	defer f2.Close()

	f3, err = os.OpenFile(appendFile+".0", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0755)
	if err != nil {
		goo_log.Error(err)
		return
	}
	defer f3.Close()

	f4, err = os.OpenFile(reduceFile+".0", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0755)
	if err != nil {
		goo_log.Error(err)
		return
	}
	defer f4.Close()

	r1 = bufio.NewReader(f1)
	r2 = bufio.NewReader(f2)

	var (
		s1, s2     string
		end1, end2 bool
	)

	for {
		if s1 == s2 {
			s1, err = r1.ReadString('\n')
			if err != nil {
				if io.EOF == err {
					err = nil
					end1 = true
				} else {
					goo_log.Error(err)
					return
				}
			}

			s2, err = r2.ReadString('\n')
			if err != nil {
				if io.EOF == err {
					err = nil
					end2 = true
				} else {
					goo_log.Error(err)
					return
				}
			}

			if end1 || end2 {
				break
			}

			continue
		}

		if s1 > s2 {
			f4.WriteString(s2)

			s2, err = r2.ReadString('\n')
			if err != nil {
				if io.EOF == err {
					err = nil
					end2 = true
				} else {
					goo_log.Error(err)
					return
				}
			}

			if end2 {
				break
			}

			continue
		}

		if s1 < s2 {
			f3.WriteString(s1)

			s1, err = r1.ReadString('\n')
			if err != nil {
				if io.EOF == err {
					err = nil
					end1 = true
				} else {
					goo_log.Error(err)
					return
				}
			}

			if end1 {
				break
			}

			continue
		}
	}

	if end1 && end2 {
		return
	}

	if end1 {
		for {
			f4.WriteString(s2)

			s2, err = r2.ReadString('\n')
			if err != nil {
				if io.EOF == err {
					err = nil
					break
				}

				goo_log.Error(err)
				return
			}
		}

		return
	}

	if end2 {
		for {
			f3.WriteString(s1)

			s1, err = r1.ReadString('\n')
			if err != nil {
				if io.EOF == err {
					err = nil
					break
				}

				goo_log.Error(err)
				return
			}
		}

		return
	}

	return
}
