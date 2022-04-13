package goo_upload

import (
	"github.com/liqiongtao/googo.io/goo"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
	"io/ioutil"
	"os"
	"path"
	"regexp"
)

type gooLocal struct {
	uploadDir string
	perm      os.FileMode
}

func Local(uploadDir string) *gooLocal {
	matched, _ := regexp.MatchString("/$", uploadDir)
	if !matched {
		uploadDir += "/"
	}
	return &gooLocal{
		uploadDir: uploadDir,
		perm:      0755,
	}
}

func (l *gooLocal) Upload(c *goo.Context) (string, error) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		return "", err
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	md5str := goo_utils.MD5(data)
	fpath := md5str[0:2] + "/" + md5str[2:4] + "/"

	if err := os.MkdirAll(l.uploadDir+fpath, l.perm); err != nil {
		return "", err
	}

	fname := fpath + md5str[8:24] + path.Ext(header.Filename)

	fw, err := os.Create(l.uploadDir + fname)
	if err != nil {
		return "", err
	}
	defer fw.Close()

	if _, err := fw.Write(data); err != nil {
		return "", err
	}

	return fname, nil
}
