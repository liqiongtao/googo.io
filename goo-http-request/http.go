package goo_http_request

import (
	"fmt"
	"io"
	"net/http"
)

func New(opts ...Option) *Request {
	r := &Request{
		Headers: map[string]string{
			"Content-Type": CONTENT_TYPE_FORM,
		},
	}
	for _, opt := range opts {
		switch opt.Name {
		case "tls":
			v := opt.Value.(map[string]string)
			r.Tls = &Tls{
				CaCrtFile:     v["caCrtFile"],
				ClientCrtFile: v["clientCrtFile"],
				ClientKeyFile: v["clientKeyFile"],
			}
		case "content-type-xml", "content-type-json", "content-type-form":
			r.SetHeader("Content-Type", opt.Value.(string))
		case "header":
			v := opt.Value.(map[string]string)
			for field, value := range v {
				r.SetHeader(field, value)
			}
		}
	}
	return r
}

func Get(url string) ([]byte, error) {
	return New().Get(url)
}

func GetWithQuery(url string, data []byte) ([]byte, error) {
	return New().GetWithQuery(url, data)
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

func GPTStream(url string, data []byte, cb func(b []byte)) error {
	return New().GPTStream(url, data, cb)
}

func Upload(url, fileField, fileName string, fh io.Reader, data map[string]string) (b []byte, err error) {
	return New().Upload(url, fileField, fileName, fh, data)
}

func Download(url, filename string) error {
	return New().Download(url, filename)
}

func SetHeader(name, value string) *Request {
	return New().SetHeader(name, value)
}

func Exists(url string) (bool, error) {
	resp, err := http.Head(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, nil
	} else if resp.StatusCode == http.StatusNotFound {
		return false, nil
	} else {
		return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
}

func Debug() *Request {
	return New().Debug()
}
