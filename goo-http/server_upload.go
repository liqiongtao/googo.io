package goo_http

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/gin-gonic/gin"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
)

type LocalUpload struct {
}

func (lu LocalUpload) Upload(c *gin.Context, uploadDir string) *Response {
	f, fh, err := c.Request.FormFile("file")
	if err != nil {
		return Error(7001, fmt.Sprintf("上传失败，原因：%s", err.Error()))
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return Error(7002, fmt.Sprintf("上传失败，原因：%s", err.Error()))
	}

	md5str := goo_utils.MD5(data)
	relDir := path.Join(md5str[0:2], md5str[2:4])
	relFile := path.Join(relDir, path.Base(fh.Filename)+"_"+md5str[8:16]+path.Ext(fh.Filename))

	if err := os.MkdirAll(path.Join(uploadDir, relDir), 0755); err != nil {
		return Error(7003, fmt.Sprintf("上传失败，原因：%s", err.Error()))
	}

	ff, err := os.Create(path.Join(uploadDir, relFile))
	if err != nil {
		return Error(7004, fmt.Sprintf("上传失败，原因：%s", err.Error()))
	}
	defer ff.Close()

	if _, err := ff.Write(data); err != nil {
		return Error(7005, fmt.Sprintf("上传失败，原因：%s", err.Error()))
	}

	return Success(gin.H{
		"url": relFile,
	})
}
