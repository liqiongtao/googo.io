package goo_request

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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
	r.client = nil
	return r
}

func (r *Request) getClient() (*http.Client, error) {
	if r.client != nil {
		return r.client, nil
	}

	timeout := r.timeout
	if timeout == 0 {
		timeout = 120 * time.Second
	}

	// 基于总超时时间动态计算 Transport 内部的各个超时阶段，确保逻辑一致性
	// 分配策略：握手和建连占用较少比例，响应等待占用较多比例
	dialTimeout := timeout / 5 // 20% 用于建立连接（DNS + TCP）
	if dialTimeout < 5*time.Second {
		dialTimeout = 5 * time.Second
	}

	tlsHandshakeTimeout := timeout / 5 // 20% 用于 TLS 握手
	if tlsHandshakeTimeout < 5*time.Second {
		tlsHandshakeTimeout = 5 * time.Second
	}

	responseHeaderTimeout := timeout - dialTimeout - tlsHandshakeTimeout
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
		ca, err := r.Tls.CaCrt()
		if err != nil {
			return nil, err
		}
		pool := x509.NewCertPool()
		if !pool.AppendCertsFromPEM(ca) {
			return nil, fmt.Errorf("failed to parse CA certificate")
		}
		cert, err := r.Tls.ClientCrt()
		if err != nil {
			return nil, err
		}
		cfg := &tls.Config{RootCAs: pool}
		if len(cert.Certificate) > 0 {
			cfg.Certificates = []tls.Certificate{cert}
		}
		transport.TLSClientConfig = cfg
	} else {
		transport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	client := &http.Client{
		Timeout: timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Transport: transport,
	}

	r.timeout = timeout
	r.client = client
	return client, nil
}

func (r *Request) Do(method, url string, reader io.Reader) (rst []byte, err error) {
	return r.do(method, url, reader, nil)
}

func (r *Request) do(method, url string, reader io.Reader, extraHeaders map[string]string) (rst []byte, err error) {
	var (
		req *http.Request
		rsp *http.Response
	)

	req, err = http.NewRequest(method, url, reader)
	if err != nil {
		if c, ok := reader.(io.Closer); ok {
			_ = c.Close()
		}
		return
	}
	defer func() {
		if req.Body != nil {
			_ = req.Body.Close()
		}
	}()

	for k, v := range r.Headers {
		req.Header.Set(k, v)
	}
	for k, v := range extraHeaders {
		req.Header.Set(k, v)
	}

	client, err := r.getClient()
	if err != nil {
		return
	}
	rsp, err = client.Do(req)
	if err != nil {
		return
	}
	defer func() { _ = rsp.Body.Close() }()

	var buf bytes.Buffer
	if _, err = io.Copy(&buf, rsp.Body); err != nil {
		return
	}

	rst = buf.Bytes()
	if rsp.StatusCode < 200 || rsp.StatusCode >= 300 {
		err = fmt.Errorf("request failed, status code: %d", rsp.StatusCode)
		return
	}

	return
}

func (r *Request) handle(method, url string, data []byte) (rsp []byte, err error) {
	return r.handleWithHeaders(method, url, data, nil)
}

func (r *Request) handleWithHeaders(method, url string, data []byte, extraHeaders map[string]string) (rsp []byte, err error) {
	rsp, err = r.do(method, url, bytes.NewReader(data), extraHeaders)
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
	if len(data) > 0 {
		if strings.Contains(url, "?") {
			url += "&" + string(data)
		} else {
			url += "?" + string(data)
		}
	}
	return r.handle("GET", url, nil)
}

func (r *Request) Post(url string, data []byte) ([]byte, error) {
	return r.handle("POST", url, data)
}

func (r *Request) PostJson(url string, data []byte) ([]byte, error) {
	// 仅对本次请求设置 Content-Type，避免污染共享 Headers
	return r.handleWithHeaders("POST", url, data, map[string]string{
		"Content-Type": CONTENT_TYPE_JSON,
	})
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
	// 流式接口按 JSON 发 body；仅改本次请求，不污染共享 Headers
	req.Header.Set("Content-Type", CONTENT_TYPE_JSON)

	base, err := r.getClient()
	if err != nil {
		return err
	}
	// Client.Timeout 覆盖整个 body 读取，长 SSE 会断；流式请求关掉总超时
	client := *base
	client.Timeout = 0

	rsp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer rsp.Body.Close()

	if rsp.StatusCode < 200 || rsp.StatusCode >= 300 {
		return fmt.Errorf("request failed, status code: %d", rsp.StatusCode)
	}

	var (
		reader   = bufio.NewReader(rsp.Body)
		headData = []byte("data: ")
		done     = "[DONE]"
	)

	for {
		b, err := reader.ReadBytes('\n')
		if len(b) > 0 {
			b2 := bytes.TrimSpace(b)
			if bytes.HasPrefix(b2, headData) {
				if cb != nil {
					// ReadBytes 复用内部缓冲，回调侧必须拷贝
					out := append([]byte(nil), b...)
					out = append(out, '\n')
					cb(out)
				}
				b3 := bytes.TrimPrefix(b2, headData)
				if string(b3) == done {
					return nil
				}
			}
		}
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
}

func (r *Request) Upload(url, fileField, fileName string, fh io.Reader, data map[string]string) ([]byte, error) {
	pr, pw := io.Pipe()
	w := multipart.NewWriter(pw)
	contentType := w.FormDataContentType()

	goo_utils.AsyncFunc(func() {
		defer pw.Close()
		defer w.Close()

		for k, v := range data {
			if err := w.WriteField(k, v); err != nil {
				_ = pw.CloseWithError(err)
				return
			}
		}

		part, err := w.CreateFormFile(fileField, fileName)
		if err != nil {
			_ = pw.CloseWithError(err)
			return
		}

		if _, err = io.Copy(part, fh); err != nil {
			_ = pw.CloseWithError(err)
			return
		}
	})

	// 仅对本次请求设置 multipart Content-Type，避免污染共享 Headers
	return r.do("POST", url, pr, map[string]string{
		"Content-Type": contentType,
	})
}

func (r *Request) Download(url, filename string) (err error) {
	filename = filepath.Clean(filename)
	if filename == "" || filename == "." || strings.HasPrefix(filename, ".."+string(filepath.Separator)) || filename == ".." {
		return fmt.Errorf("invalid download filename")
	}

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
		dirname := filepath.Dir(filename)
		if dirname != "" && dirname != "." {
			_ = os.MkdirAll(dirname, 0755)
		}
	}

	var f *os.File
	tmpName := filename + ".0"
	ok := false
	{
		f, err = os.Create(tmpName)
		if err != nil {
			return
		}
		defer func() {
			_ = f.Close()
			if !ok {
				_ = os.Remove(tmpName)
			}
		}()
	}

	var req *http.Request
	{
		req, err = http.NewRequest("GET", url, nil)
		if err != nil {
			return
		}
		for k, v := range r.Headers {
			req.Header.Set(k, v)
		}
	}

	savedTimeout, savedClient := r.timeout, r.client
	if r.timeout == 0 {
		r.timeout = 5 * time.Minute
		r.client = nil
	}
	defer func() {
		r.timeout = savedTimeout
		r.client = savedClient
	}()

	var resp *http.Response
	{
		var client *http.Client
		client, err = r.getClient()
		if err != nil {
			return
		}
		resp, err = client.Do(req)
		if err != nil {
			return
		}
		defer resp.Body.Close()
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		err = fmt.Errorf("download failed, status code: %d", resp.StatusCode)
		return
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

	if err = os.Rename(tmpName, filename); err != nil {
		return
	}
	ok = true
	return nil
}
