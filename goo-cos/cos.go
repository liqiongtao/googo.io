package goo_cos

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	goo_file "github.com/liqiongtao/googo.io/goo-file"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"github.com/tencentyun/cos-go-sdk-v5"
)

type CosClient struct {
	*cos.Client
	Config CosConfig
}

func NewCosClient(cfg CosConfig) (*CosClient, error) {
	bucketUrl, err := url.Parse(cfg.BucketURL())
	if err != nil {
		return nil, fmt.Errorf("invalid cos bucket url: %w", err)
	}

	c := cos.NewClient(
		&cos.BaseURL{BucketURL: bucketUrl},
		&http.Client{
			Transport: &cos.AuthorizationTransport{
				SecretID:     cfg.SecretId,
				SecretKey:    cfg.SecretKey,
				SessionToken: cfg.SessionToken,
				Transport: &http.Transport{
					// 1. 连接复用与超时（核心，防泄露）
					IdleConnTimeout:       90 * time.Second, // 空闲连接超时：高并发下可适当延长（30s→60s），提升复用率
					ResponseHeaderTimeout: 15 * time.Second, // 响应头超时：高并发下服务端可能慢，适度放宽（10s→15s）
					TLSHandshakeTimeout:   10 * time.Second, // TLS 握手超时：高并发下握手可能排队，放宽（5s→10s）

					// 2. 连接建立超时
					DialContext: (&net.Dialer{
						Timeout:   10 * time.Second, // 拨号超时：高并发下网络可能拥塞，放宽（5s→10s）
						KeepAlive: 60 * time.Second, // TCP 保活：保持长连接，提升复用
					}).DialContext,

					// 3. 连接池大小（高并发核心调优）
					MaxIdleConns:        1000,  // 全局最大空闲连接：默认100，高并发下需大幅提升（根据QPS调整）
					MaxIdleConnsPerHost: 200,   // 单Host最大空闲连接：默认2，高并发下必须调高（比如COS域名）
					MaxConnsPerHost:     500,   // 单Host最大并发连接：默认无限制，限制避免压垮服务端
					DisableCompression:  false, // 启用压缩：减少传输量，提升高并发下的吞吐量

					// 4. 其他高并发优化
					ExpectContinueTimeout: 2 * time.Second, // 处理 Expect: 100-Continue 的超时，缩短等待
				},
			},
		},
	)

	return &CosClient{
		Client: c,
		Config: cfg,
	}, nil
}

func trimObjectKey(objectKey string) string {
	return strings.TrimPrefix(objectKey, "/")
}

func closeCOSResponse(rsp *cos.Response) {
	if rsp != nil && rsp.Body != nil {
		_ = rsp.Body.Close()
	}
}

// 上传文件
func (c *CosClient) Upload(localFileName, objectKey string) error {
	objectKey = trimObjectKey(objectKey)

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
	defer closeCOSResponse(rsp)

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
	objectKey = trimObjectKey(objectKey)

	if err := os.MkdirAll(path.Dir(localFileName), 0755); err != nil {
		goo_log.ErrorF("create dir %s error: %s", path.Dir(localFileName), err.Error())
		return err
	}

	rsp, err := c.Object.GetToFile(context.TODO(), objectKey, localFileName, nil)
	defer closeCOSResponse(rsp)
	if err != nil {
		goo_log.ErrorF("download %s error: %s", objectKey, err.Error())
		return err
	}

	if !goo_file.Exist(localFileName) {
		goo_log.ErrorF("download %s error", objectKey)
		return fmt.Errorf("下载 %s 失败", objectKey)
	}

	return nil
}

// 获取文件内容
func (c *CosClient) Get(objectKey string) ([]byte, error) {
	objectKey = trimObjectKey(objectKey)

	// 使用cos客户端获取文件内容
	resp, err := c.Object.Get(context.Background(), objectKey, nil)
	if err != nil {
		closeCOSResponse(resp)
		goo_log.ErrorF("get %s error: %s", objectKey, err.Error())
		return nil, err
	}
	defer closeCOSResponse(resp)

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

		res, rsp, err := c.Bucket.Get(context.TODO(), opt)
		if err != nil {
			closeCOSResponse(rsp)
			return err
		}
		closeCOSResponse(rsp)

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
	objectKey = trimObjectKey(objectKey)

	rsp, err := c.Object.Head(context.Background(), objectKey, nil)
	if err != nil {
		closeCOSResponse(rsp)
		goo_log.ErrorF("head %s error: %s", objectKey, err.Error())
		return nil, err
	}
	// Head 一般无 body，仍安全关闭，避免连接泄漏
	closeCOSResponse(rsp)

	return rsp, nil
}

func (c *CosClient) IsExist(objectKey string) bool {
	objectKey = trimObjectKey(objectKey)

	ok, err := c.Object.IsExist(context.Background(), objectKey)
	if err != nil {
		goo_log.ErrorF("exist %s error: %s", objectKey, err.Error())
		return false
	}

	return ok
}

func (c *CosClient) Delete(objectKey string) error {
	objectKey = trimObjectKey(objectKey)

	rsp, err := c.Object.Delete(context.Background(), objectKey, nil)
	defer closeCOSResponse(rsp)
	if err != nil {
		goo_log.ErrorF("delete %s error: %s", objectKey, err.Error())
		return err
	}

	return nil
}

func (c *CosClient) Copy(sourceObjectKey, targetObjectKey string, targetCosClient *cos.Client) error {
	sourceObjectKey = trimObjectKey(sourceObjectKey)
	targetObjectKey = trimObjectKey(targetObjectKey)

	if !strings.HasPrefix(sourceObjectKey, c.BaseURL.BucketURL.Host+c.BaseURL.BucketURL.Path) {
		sourceObjectKey = c.BaseURL.BucketURL.Host + c.BaseURL.BucketURL.Path + "/" + sourceObjectKey
	}

	_, rsp, err := targetCosClient.Object.Copy(context.Background(), targetObjectKey, sourceObjectKey, nil)
	defer closeCOSResponse(rsp)
	if err != nil {
		goo_log.ErrorF("copy %s error: %s", sourceObjectKey, err.Error())
		return err
	}

	return nil
}
