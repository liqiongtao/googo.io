package goofile

import (
	"bufio"
	"errors"
	"io"
	"os"

	goolog "github.com/liqiongtao/googo.io/goo-log"
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

		if e := os.Rename(appendFile+".0", appendFile); e != nil {
			err = e
			goolog.Error(e)
			return
		}
		if e := os.Rename(reduceFile+".0", reduceFile); e != nil {
			err = e
			goolog.Error(e)
		}
	}()

	var (
		f1, f2, f3, f4 *os.File
		r1, r2         *bufio.Reader
	)

	f1, err = os.OpenFile(srcFile, os.O_RDONLY, 0)
	if err != nil {
		goolog.Error(err)
		return
	}
	defer f1.Close()

	f2, err = os.OpenFile(targetFile, os.O_RDONLY, 0)
	if err != nil {
		goolog.Error(err)
		return
	}
	defer f2.Close()

	f3, err = os.OpenFile(appendFile+".0", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0755)
	if err != nil {
		goolog.Error(err)
		return
	}
	defer f3.Close()

	f4, err = os.OpenFile(reduceFile+".0", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0755)
	if err != nil {
		goolog.Error(err)
		return
	}
	defer f4.Close()

	r1 = bufio.NewReader(f1)
	r2 = bufio.NewReader(f2)

	write := func(f *os.File, s string) bool {
		if _, e := f.WriteString(s); e != nil {
			err = e
			goolog.Error(e)
			return false
		}
		return true
	}

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
					goolog.Error(err)
					return
				}
			}

			s2, err = r2.ReadString('\n')
			if err != nil {
				if io.EOF == err {
					err = nil
					end2 = true
				} else {
					goolog.Error(err)
					return
				}
			}

			if end1 || end2 {
				break
			}

			continue
		}

		if s1 > s2 {
			if !write(f4, s2) {
				return
			}

			s2, err = r2.ReadString('\n')
			if err != nil {
				if io.EOF == err {
					err = nil
					end2 = true
				} else {
					goolog.Error(err)
					return
				}
			}

			if end2 {
				break
			}

			continue
		}

		if s1 < s2 {
			if !write(f3, s1) {
				return
			}

			s1, err = r1.ReadString('\n')
			if err != nil {
				if io.EOF == err {
					err = nil
					end1 = true
				} else {
					goolog.Error(err)
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
		if s1 != "" && s2 != "" && s1 != s2 {
			if s1 < s2 {
				if !write(f3, s1) || !write(f4, s2) {
					return
				}
			} else {
				if !write(f4, s2) || !write(f3, s1) {
					return
				}
			}
		} else if s1 != "" && s2 == "" {
			if !write(f3, s1) {
				return
			}
		} else if s2 != "" && s1 == "" {
			if !write(f4, s2) {
				return
			}
		}
		return
	}

	if end1 {
		// 先结算残留的 s1/s2，再排空 file2 剩余行
		if s1 != "" && s2 != "" {
			if s1 != s2 {
				if s1 < s2 {
					if !write(f3, s1) || !write(f4, s2) {
						return
					}
				} else {
					if !write(f4, s2) || !write(f3, s1) {
						return
					}
				}
			}
		} else if s1 != "" {
			if !write(f3, s1) {
				return
			}
		} else if s2 != "" {
			if !write(f4, s2) {
				return
			}
		}

		for {
			s2, err = r2.ReadString('\n')
			if err != nil {
				if io.EOF == err {
					if s2 != "" {
						if !write(f4, s2) {
							return
						}
					}
					err = nil
					break
				}
				goolog.Error(err)
				return
			}
			if !write(f4, s2) {
				return
			}
		}
		return
	}

	if end2 {
		if s1 != "" && s2 != "" {
			if s1 != s2 {
				if s1 < s2 {
					if !write(f3, s1) || !write(f4, s2) {
						return
					}
				} else {
					if !write(f4, s2) || !write(f3, s1) {
						return
					}
				}
			}
		} else if s2 != "" {
			if !write(f4, s2) {
				return
			}
		} else if s1 != "" {
			if !write(f3, s1) {
				return
			}
		}

		for {
			s1, err = r1.ReadString('\n')
			if err != nil {
				if io.EOF == err {
					if s1 != "" {
						if !write(f3, s1) {
							return
						}
					}
					err = nil
					break
				}
				goolog.Error(err)
				return
			}
			if !write(f3, s1) {
				return
			}
		}
		return
	}

	return
}
