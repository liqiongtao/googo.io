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
	client  *http.Client
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
	if r.timeout.Seconds() == 0 {
		r.timeout = 120 * time.Second
	}

	if r.client != nil {
		return r.client
	}

	// 基于总超时时间动态计算 Transport 内部的各个超时阶段，确保逻辑一致性
	// 分配策略：握手和建连占用较少比例，响应等待占用较多比例
	dialTimeout := r.timeout / 5 // 20% 用于建立连接（DNS + TCP）
	if dialTimeout < 5*time.Second {
		dialTimeout = 5 * time.Second
	}

	tlsHandshakeTimeout := r.timeout / 5 // 20% 用于 TLS 握手
	if tlsHandshakeTimeout < 5*time.Second {
		tlsHandshakeTimeout = 5 * time.Second
	}

	responseHeaderTimeout := r.timeout - dialTimeout - tlsHandshakeTimeout
	if responseHeaderTimeout < 10*time.Second {
		responseHeaderTimeout = 10 * time.Second
	}

	transport := &http.Transport{
		// 1. 连接复用与超时
		IdleConnTimeout:       90 * time.Second, // 空闲连接保持时间，独立于请求超时
		ResponseHeaderTimeout: responseHeaderTimeout,
		TLSHandshakeTimeout:   tlsHandshakeTimeout,

		// 2. 连接建立超时
		DialContext: (&net.Dialer{
			Timeout:   dialTimeout,
			KeepAlive: 60 * time.Second,
		}).DialContext,

		// 3. 连接池大小（高并发核心调优）
		MaxIdleConns:        1000,
		MaxIdleConnsPerHost: 200,
		MaxConnsPerHost:     500,
		DisableCompression:  false,

		// 4. 其他优化
		ExpectContinueTimeout: 2 * time.Second,
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

	client := &http.Client{
		Timeout: r.timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Transport: transport,
	}

	r.client = client
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

	var buf bytes.Buffer
	if _, err = io.Copy(&buf, rsp.Body); err != nil {
		return
	}

	rst = buf.Bytes()

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
