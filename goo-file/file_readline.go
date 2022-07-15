package goo_file

import (
	"bufio"
	"errors"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"io"
	"os"
)

func ReadByLine(filename string, cb func(b []byte, end bool) error) error {
	if !Exist(filename) {
		return errors.New("文件不存在")
	}

	f, err := os.OpenFile(filename, os.O_RDWR, 0755)
	if err != nil {
		goo_log.Error(err)
		return err
	}
	defer f.Close()

	r := bufio.NewReader(f)
	for {
		b, err := r.ReadBytes('\n')

		if err != nil {
			if io.EOF == err {
				return cb([]byte{}, true)
			}

			goo_log.Error(err)
			return err
		}

		if err := cb(b, false); err != nil {
			return err
		}
	}
}
