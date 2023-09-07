package goo_file

import (
	"crypto/md5"
	"encoding/hex"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"io"
	"io/ioutil"
	"os"
)

func MD5(file string) (string, error) {
	defaultSize := int64(16 * 1024 * 1024)

	info, err := os.Stat(file)
	if err != nil {
		goo_log.Error(err)
		return "", err
	}

	// 小文件
	if info.Size() < defaultSize {
		b, _ := ioutil.ReadFile(file)
		sum := md5.Sum(b)
		return hex.EncodeToString(sum[:]), nil
	}

	// 大文件
	{
		tempFile, err := ioutil.TempFile(os.TempDir(), "goo-md5-temp-file")
		if err != nil {
			goo_log.Error(err)
			return "", err
		}
		defer tempFile.Close()

		f, err := os.OpenFile(file, os.O_RDONLY, 0755)
		if err != nil {
			goo_log.Error(err)
			return "", err
		}
		defer f.Close()

		io.Copy(tempFile, f)
		tempFile.Seek(0, os.SEEK_SET)

		h := md5.New()
		io.Copy(h, tempFile)
		return hex.EncodeToString(h.Sum(nil)), nil
	}
}
