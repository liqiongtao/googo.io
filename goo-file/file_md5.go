package goo_file

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"

	goo_log "github.com/liqiongtao/googo.io/goo-log"
)

func MD5(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		goo_log.Error(err)
		return "", err
	}
	defer f.Close()

	h := md5.New()
	if _, err = io.Copy(h, f); err != nil {
		goo_log.Error(err)
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
