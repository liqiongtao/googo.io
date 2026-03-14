package goo_http_request

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"path"
	"time"

	goo_log "github.com/liqiongtao/googo.io/goo-log"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
)

type Request struct {
	Headers map[string]string
	Tls     *Tls
	timeout time.Duration
	debug   bool
}

func (r *Request) Debug() *Request {
	r.debug = true
	return r
}

func (r *Request) SetHeader(name, value string) *Request {
	r.Headers[name] = value
	return r
}

func (r *Request) SetContentType(contentType string) *Request {
	r.SetHeader("Content-Type", contentType)
	return r
}

func (r *Request) JsonContentType() *Request {
	r.SetHeader("Content-Type", CONTENT_TYPE_JSON)
	return r
}

func (r *Request) SetTimeout(d time.Duration) *Request {
	r.timeout = d
	return r
}

func (r *Request) getClient() *http.Client {
	if r.timeout == 0 {
		r.timeout = 30 * time.Second
	}
	client := &http.Client{
		Timeout: r.timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	transport := &http.Transport{
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
	}

	if r.Tls != nil {
		pool := x509.NewCertPool()
		pool.AppendCertsFromPEM(r.Tls.CaCrt())
		transport.TLSClientConfig = &tls.Config{
			RootCAs:      pool,
			Certificates: []tls.Certificate{r.Tls.ClientCrt()},
		}
	} else {
		transport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	client.Transport = transport

	return client
}

func (r *Request) Do(method, url string, reader io.Reader) (rst []byte, err error) {
	var (
		req *http.Request
		rsp *http.Response
	)

	req, err = http.NewRequest(method, url, reader)
	if err != nil {
		return
	}
	defer func() { _ = req.Body.Close() }()

	for k, v := range r.Headers {
		req.Header.Set(k, v)
	}

	rsp, err = r.getClient().Do(req)
	if err != nil {
		return
	}
	defer func() { _ = rsp.Body.Close() }()

	var (
		bf bytes.Buffer
		n  int
	)

	for {
		bts := make([]byte, 1024)
		n, err = rsp.Body.Read(bts)
		if err != nil && err != io.EOF {
			return
		}
		if n == 0 {
			err = nil
			break
		}

		bf.Write(bts[:n])
	}

	rst = bf.Bytes()

	return
}

func (r *Request) handle(method, url string, data []byte) (rsp []byte, err error) {
	rsp, err = r.Do(method, url, bytes.NewReader(data))
	if r.debug {
		l := goo_log.WithTag(TAG).
			WithField("method", method).
			WithField("url", url).
			WithField("header", r.Headers).
			WithField("request-data", string(data)).
			WithField("response", string(rsp))
		if err != nil && err != io.EOF {
			l.Error(err)
		} else {
			l.Debug()
		}
	}
	return
}

func (r *Request) Get(url string) ([]byte, error) {
	return r.handle("GET", url, nil)
}

func (r *Request) GetWithQuery(url string, data []byte) ([]byte, error) {
	return r.handle("GET", url, data)
}

func (r *Request) Post(url string, data []byte) ([]byte, error) {
	return r.handle("POST", url, data)
}

func (r *Request) PostJson(url string, data []byte) ([]byte, error) {
	return r.JsonContentType().handle("POST", url, data)
}

func (r *Request) Put(url string, data []byte) ([]byte, error) {
	return r.handle("PUT", url, data)
}

func (r *Request) GPTStream(url string, data []byte, cb func(b []byte)) error {
	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return err
	}

	for k, v := range r.Headers {
		req.Header.Set(k, v)
	}

	rsp, err := r.getClient().Do(req)
	if err != nil {
		return err
	}

	defer rsp.Body.Close()

	var (
		reader   = bufio.NewReader(rsp.Body)
		headData = []byte("data: ")
		done     = "[DONE]"
	)

	for {
		b, err := reader.ReadBytes('\n')
		if err != nil {
			return err
		}

		b2 := bytes.TrimSpace(b)
		if !bytes.HasPrefix(b2, headData) {
			continue
		}

		cb(append(b, '\n'))

		b3 := bytes.TrimPrefix(b2, headData)
		if string(b3) == done {
			break
		}
	}

	return nil
}

func (r *Request) Upload(url, fileField, fileName string, fh io.Reader, data map[string]string) ([]byte, error) {
	pr, pw := io.Pipe()
	w := multipart.NewWriter(pw)

	goo_utils.AsyncFunc(func() {
		for k, v := range data {
			w.WriteField(k, v)
		}

		part, _ := w.CreateFormFile(fileField, fileName)

		io.CopyBuffer(part, fh, nil)

		w.Close()
		pw.Close()
	})

	r.SetHeader("Content-Type", w.FormDataContentType())

	return r.Do("POST", url, pr)
}

func (r *Request) Download(url, filename string) (err error) {
	defer func() {
		if !r.debug {
			return
		}
		l := goo_log.WithTag(TAG).WithField("url", url).WithField("file", filename)
		if err != nil {
			l.Error("下载失败", err)
			return
		}
		l.Debug("下载成功")
	}()

	// 创建目录
	{
		dirname := path.Dir(filename)
		if dirname != "" && dirname != "." && dirname != "./" {
			os.MkdirAll(dirname, 0755)
		}
	}

	var f *os.File
	{
		f, err = os.Create(filename + ".0")
		if err != nil {
			return
		}
		defer f.Close()
	}

	var req *http.Request
	{
		req, err = http.NewRequest("GET", url, nil)
		if err != nil {
			return
		}
	}

	if r.timeout == 0 {
		r.timeout = 5 * time.Minute
	}

	var resp *http.Response
	{
		resp, err = r.getClient().Do(req)
		if err != nil {
			return
		}
		defer resp.Body.Close()
	}

	for {
		var (
			n   int
			bts = make([]byte, 1024)
		)

		n, err = resp.Body.Read(bts)
		if err != nil && err != io.EOF {
			return
		}
		if n == 0 {
			err = nil
			break
		}

		if _, err = f.Write(bts[:n]); err != nil {
			return
		}
	}

	return os.Rename(filename+".0", filename)
}
