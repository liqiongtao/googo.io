package goo_http_request

import (
	"bytes"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"io"
	"mime/multipart"
)

func New(opts ...Option) *Request {
	req := &Request{
		Headers: map[string]string{
			"Content-Type": CONTENT_TYPE_FORM,
		},
	}
	for _, opt := range opts {
		switch opt.Name {
		case "tsl":
			v := opt.Value.(map[string]string)
			req.Tls = &Tls{
				CaCrtFile:     v["caCrtFile"],
				ClientCrtFile: v["clientCrtFile"],
				ClientKeyFile: v["clientKeyFile"],
			}
		case "content-type-xml", "content-type-json", "content-type-form":
			req.Headers["Content-Type"] = opt.Value.(string)
		case "header":
			v := opt.Value.(map[string]string)
			for field, value := range v {
				req.Headers[field] = value
			}
		}
	}
	return req
}

func Get(url string) ([]byte, error) {
	return New().Get(url)
}

func Post(url string, data []byte) ([]byte, error) {
	return New().Post(url, data)
}

func PostJson(url string, data []byte) ([]byte, error) {
	return New().JsonContentType().Post(url, data)
}

func Put(url string, data []byte) ([]byte, error) {
	return New().Put(url, data)
}

func Upload(url, field, fileName string, f io.Reader, data map[string]string) ([]byte, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(field, fileName)
	if err != nil {
		goo_log.Error(err.Error())
		return nil, err
	}
	if _, err = io.Copy(part, f); err != nil {
		goo_log.Error(err.Error())
		return nil, err
	}

	for k, v := range data {
		writer.WriteField(k, v)
	}

	if err = writer.Close(); err != nil {
		goo_log.Error(err.Error())
		return nil, err
	}

	request := New()
	request.SetHearder("Content-Type", writer.FormDataContentType())
	return request.Do("POST", url, body)
}
