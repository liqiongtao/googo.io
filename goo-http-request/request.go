package goo_http_request

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"time"
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
		r.timeout = 8 * time.Second
	}
	client := &http.Client{
		Timeout: r.timeout,
	}
	if r.Tls != nil {
		pool := x509.NewCertPool()
		pool.AppendCertsFromPEM(r.Tls.CaCrt())
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      pool,
				Certificates: []tls.Certificate{r.Tls.ClientCrt()},
			},
		}
	}
	return client
}

func (r *Request) Do(method, url string, body []byte) ([]byte, error) {
	l := goo_log.WithTag("goo-http-request")

	if r.debug {
		l.WithField("method", method)
		l.WithField("url", url)
		l.WithField("data", string(body))
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		l.Error(err)
		return nil, err
	}

	for k, v := range r.Headers {
		req.Header.Set(k, v)
	}

	if r.debug {
		l.WithField("header", r.Headers)
	}

	rsp, err := r.getClient().Do(req)
	if err != nil {
		l.Error(err)
		return nil, err
	}

	defer rsp.Body.Close()

	buf, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		l.Error(err)
		return nil, err
	}

	if r.debug {
		l.WithField("response", string(buf))
		l.Debug()
	}

	return buf, nil
}

func (r *Request) Get(url string) ([]byte, error) {
	return r.Do("GET", url, nil)
}

func (r *Request) GetWithQuery(url string, data []byte) ([]byte, error) {
	return r.Do("GET", url, data)
}

func (r *Request) Post(url string, data []byte) ([]byte, error) {
	return r.Do("POST", url, data)
}

func (r *Request) PostJson(url string, data []byte) ([]byte, error) {
	return r.JsonContentType().Do("POST", url, data)
}

func (r *Request) Put(url string, data []byte) ([]byte, error) {
	return r.Do("PUT", url, data)
}

func (r *Request) Upload(url, fileField, fileName string, f io.Reader, data map[string]string) (b []byte, err error) {
	var (
		body bytes.Buffer
		part io.Writer
	)

	w := multipart.NewWriter(&body)
	if part, err = w.CreateFormFile(fileField, fileName); err != nil {
		goo_log.WithTag("[http-request-upload]").Error(err.Error())
		return
	}

	if _, err = io.Copy(part, f); err != nil {
		goo_log.WithTag("[http-request-upload]").Error(err.Error())
		return
	}

	for k, v := range data {
		w.WriteField(k, v)
	}

	w.Close()

	r.SetHeader("Content-Type", w.FormDataContentType())
	return r.Do("POST", url, body.Bytes())
}
