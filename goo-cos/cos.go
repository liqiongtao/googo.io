package goo_cos

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	goo_file "github.com/liqiongtao/googo.io/goo-file"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"github.com/tencentyun/cos-go-sdk-v5"
)

type CosClient struct {
	*cos.Client
	Config CosConfig
}

func NewCosClient(cfg CosConfig) *CosClient {
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

	return &CosClient{
		Client: c,
		Config: cfg,
	}
}

// 上传文件
func (c *CosClient) Upload(localFileName, objectKey string) error {
	f, err := os.Open(localFileName)
	if err != nil {
		goo_log.ErrorF("open file %s error: %s", localFileName, err.Error())
		return err
	}
	defer func() { _ = f.Close() }()

	// 获取文件大小
	stat, err := f.Stat()
	if err != nil {
		goo_log.ErrorF("stat file %s error: %s", localFileName, err.Error())
		return err
	}

	// 执行上传
	rsp, err := c.Object.Put(context.Background(), objectKey, f, &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentLength: stat.Size(),
		},
	})
	defer func() { _ = rsp.Body.Close() }()

	if err != nil {
		goo_log.ErrorF("put %s error: %s", objectKey, err.Error())
		return err
	}
	if rsp != nil && rsp.StatusCode != 200 {
		goo_log.ErrorF("put %s error, status code: %d", objectKey, rsp.StatusCode)
		return fmt.Errorf("上传文件失败，状态码: %d", rsp.StatusCode)
	}

	return nil
}

// 下载文件
func (c *CosClient) Download(localFileName, objectKey string) error {
	if err := os.MkdirAll(path.Dir(localFileName), 0755); err != nil {
		goo_log.ErrorF("create dir %s error: %s", path.Dir(localFileName), err.Error())
		return err
	}

	rsp, err := c.Object.GetToFile(context.TODO(), objectKey, localFileName, nil)
	if err != nil {
		goo_log.ErrorF("download %s error: %s", objectKey, err.Error())
		return err
	}
	defer func() { _ = rsp.Body.Close() }()

	if !goo_file.Exist(localFileName) {
		goo_log.ErrorF("download %s error", objectKey)
		return fmt.Errorf("下载 %s 失败", objectKey)
	}

	return nil
}

// 获取文件内容
func (c *CosClient) Get(objectKey string) ([]byte, error) {
	if objectKey[0:1] == "/" {
		objectKey = objectKey[1:]
	}

	// 使用cos客户端获取文件内容
	resp, err := c.Object.Get(context.Background(), objectKey, nil)
	if err != nil {
		goo_log.ErrorF("get %s error: %s", objectKey, err.Error())
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	// 读取响应体内容
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		goo_log.ErrorF("read %s body error: %s", objectKey, err.Error())
		return nil, err
	}

	return b, nil
}

func (c *CosClient) List(fn func(objectKey cos.Object) error) error {
	var marker string

	opt := &cos.BucketGetOptions{
		MaxKeys: 1000, // 每次最多列出1000个对象
	}

	for {
		opt.Marker = marker

		res, _, err := c.Bucket.Get(context.TODO(), opt)
		if err != nil {
			return err
		}

		for _, v := range res.Contents {
			if err = fn(v); err != nil {
				return err
			}
		}

		if res.IsTruncated {
			marker = res.NextMarker
		} else {
			break
		}
	}

	return nil
}

func (c *CosClient) Head(objectKey string) (*cos.Response, error) {
	if objectKey[0:1] == "/" {
		objectKey = objectKey[1:]
	}

	rsp, err := c.Object.Head(context.Background(), objectKey, nil)
	if err != nil {
		goo_log.ErrorF("head %s error: %s", objectKey, err.Error())
		return nil, err
	}
	defer func() { _ = rsp.Body.Close() }()

	return rsp, nil
}

func (c *CosClient) IsExist(objectKey string) bool {
	if objectKey[0:1] == "/" {
		objectKey = objectKey[1:]
	}

	_, err := c.Object.IsExist(context.Background(), objectKey)
	if err != nil {
		goo_log.ErrorF("exist %s error: %s", objectKey, err.Error())
		return false
	}

	return true
}

func (c *CosClient) Delete(objectKey string) error {
	if objectKey[0:1] == "/" {
		objectKey = objectKey[1:]
	}

	// 检查对象是否存在
	_, err := c.Object.Delete(context.Background(), objectKey, nil)
	if err != nil {
		goo_log.ErrorF("delete %s error: %s", objectKey, err.Error())
		return err
	}

	return nil
}

func (c *CosClient) Copy(sourceObjectKey, targetObjectKey string, targetCosClient *cos.Client) error {
	if sourceObjectKey[0:1] == "/" {
		sourceObjectKey = sourceObjectKey[1:]
	}
	if !strings.HasPrefix(sourceObjectKey, c.BaseURL.BucketURL.Host+c.BaseURL.BucketURL.Path) {
		sourceObjectKey = c.BaseURL.BucketURL.Host + c.BaseURL.BucketURL.Path + "/" + sourceObjectKey
	}

	if targetObjectKey[0:1] == "/" {
		targetObjectKey = targetObjectKey[1:]
	}

	_, rsp, err := targetCosClient.Object.Copy(context.Background(), targetObjectKey, sourceObjectKey, nil)
	if err != nil {
		goo_log.ErrorF("copy %s error: %s", sourceObjectKey, err.Error())
		return err
	}
	defer func() { _ = rsp.Body.Close() }()

	return nil
}
