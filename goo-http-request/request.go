package goo_http_request

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	goo_log "googo.io/goo-log"
	"io"
	"io/ioutil"
	"net/http"
)

type Request struct {
	Headers map[string]string
	Tls     *Tls
	debug   bool
}

func (r *Request) Debug() *Request {
	r.debug = true
	return r
}

func (r *Request) SetHearder(name, value string) *Request {
	r.Headers[name] = value
	return r
}

func (r *Request) SetContentType(contentType string) *Request {
	r.SetHearder("Content-Type", contentType)
	return r
}

func (r *Request) JsonContentType() *Request {
	r.SetHearder("Content-Type", CONTENT_TYPE_JSON)
	return r
}

func (r *Request) getClient() *http.Client {
	client := &http.Client{}
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

func (r *Request) Do(method, url string, body io.Reader) ([]byte, error) {
	var l *goo_log.Entry
	if r.debug {
		l = goo_log.WithField("method", method).WithField("url", url)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		if r.debug {
			l.Debug("[http-request]")
		}
		goo_log.Error(err.Error())
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
		if r.debug {
			l.Debug("[http-request]")
		}
		goo_log.Error(err.Error())
		return nil, err
	}

	defer rsp.Body.Close()

	buf, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		if r.debug {
			l.Debug("[http-request]")
		}
		goo_log.Error(err.Error())
		return nil, err
	}

	if r.debug {
		l.WithField("response", string(buf))
		l.Debug("[http-request]")
	}

	return buf, nil
}

func (r *Request) Get(url string) ([]byte, error) {
	return r.Do("GET", url, nil)
}

func (r *Request) Post(url string, data []byte) ([]byte, error) {
	return r.Do("POST", url, bytes.NewReader(data))
}
