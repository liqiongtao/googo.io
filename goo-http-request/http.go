package goo_http_request

import (
	"io"
)

func New(opts ...Option) *Request {
	req := &Request{
		Headers: map[string]string{
			"Content-Type": CONTENT_TYPE_FORM,
		},
	}
	for _, opt := range opts {
		switch opt.Name {
		case "tls":
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

func Upload(url, fileField, fileName string, f io.Reader, data map[string]string) (b []byte, err error) {
	return New().Upload(url, fileField, fileName, f, data)
}
