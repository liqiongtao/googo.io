package goofile

import (
	"bufio"
	"errors"
	goolog "github.com/liqiongtao/googo.io/goo-log"
	"io"
	"os"
)

func ReadByLine(filename string, cb func(b []byte, end bool) error) error {
	if !Exist(filename) {
		return errors.New("文件不存在")
	}

	f, err := os.OpenFile(filename, os.O_RDONLY, 0)
	if err != nil {
		goolog.Error(err)
		return err
	}
	defer f.Close()

	r := bufio.NewReader(f)
	for {
		b, err := r.ReadBytes('\n')

		if err != nil {
			if io.EOF == err {
				if len(b) == 0 {
					return cb(nil, true)
				}
				return cb(b, true)
			}

			goolog.Error(err)
			return err
		}

		if err := cb(b, false); err != nil {
			return err
		}
	}
}
