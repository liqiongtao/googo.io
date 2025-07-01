package goo_cos

import (
	"context"
	"fmt"
	goo_file "github.com/liqiongtao/googo.io/goo-file"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"github.com/tencentyun/cos-go-sdk-v5"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
)

func CosClient(cfg CosConfig) *cos.Client {
	bucketUrl, _ := url.Parse(cfg.BucketURL())

	c := cos.NewClient(
		&cos.BaseURL{BucketURL: bucketUrl},
		&http.Client{
			Transport: &cos.AuthorizationTransport{
				SecretID:     cfg.SecretId,
				SecretKey:    cfg.SecretKey,
				SessionToken: cfg.SessionToken,
			},
		},
	)

	return c
}

// 上传文件
func CosUpload(localFileName, objectKey string, c *cos.Client) error {
	f, err := os.Open(localFileName)
	if err != nil {
		goo_log.ErrorF("open file %s error", localFileName)
		return err
	}
	defer f.Close()

	// 获取文件大小
	stat, err := f.Stat()
	if err != nil {
		goo_log.ErrorF("stat file %s error", localFileName)
		return err
	}

	// 执行上传
	resp, err := c.Object.Put(context.Background(), objectKey, f, &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentLength: stat.Size(),
		},
	})
	if err != nil {
		goo_log.ErrorF("put %s error", objectKey)
		return err
	}
	if resp != nil && resp.StatusCode != 200 {
		goo_log.ErrorF("put %s error, status code: %d", objectKey, resp.StatusCode)
		return fmt.Errorf("上传文件失败，状态码: %d", resp.StatusCode)
	}

	return nil
}

// 下载文件
func CosDownload(localFileName, objectKey string, c *cos.Client) error {
	if err := os.MkdirAll(path.Dir(localFileName), 0755); err != nil {
		goo_log.ErrorF("create dir %s error", path.Dir(localFileName))
		return err
	}

	if _, err := c.Object.GetToFile(context.TODO(), objectKey, localFileName, nil); err != nil {
		goo_log.ErrorF("download %s error", objectKey)
		return err
	}

	if !goo_file.Exist(localFileName) {
		goo_log.ErrorF("download %s error", objectKey)
		return fmt.Errorf("下载 %s 失败", objectKey)
	}

	return nil
}

// 获取文件内容
func CosGet(objectKey string, c *cos.Client) ([]byte, error) {
	if objectKey[0:1] == "/" {
		objectKey = objectKey[1:]
	}

	// 使用cos客户端获取文件内容
	resp, err := c.Object.Get(context.Background(), objectKey, nil)
	if err != nil {
		goo_log.ErrorF("get %s error", objectKey)
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应体内容
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		goo_log.ErrorF("read %s body error", objectKey)
		return nil, err
	}

	return b, nil
}

func CosExists(objectKey string, c *cos.Client) bool {
	if objectKey[0:1] == "/" {
		objectKey = objectKey[1:]
	}

	// 检查对象是否存在
	_, err := c.Object.Head(context.Background(), objectKey, nil)
	if err != nil {
		goo_log.ErrorF("head %s error", objectKey)
		return false
	}

	return true
}

func CosDelete(objectKey string, c *cos.Client) error {
	if objectKey[0:1] == "/" {
		objectKey = objectKey[1:]
	}

	// 检查对象是否存在
	_, err := c.Object.Delete(context.Background(), objectKey, nil)
	if err != nil {
		goo_log.ErrorF("delete %s error", objectKey)
		return err
	}

	return nil
}
