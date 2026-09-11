package goofile

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"

	goolog "github.com/liqiongtao/googo.io/goo-log"
)

func MD5(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		goolog.Error(err)
		return "", err
	}
	defer f.Close()

	h := md5.New()
	if _, err = io.Copy(h, f); err != nil {
		goolog.Error(err)
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
