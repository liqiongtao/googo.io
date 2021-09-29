package goo_http_request

import (
	"crypto/tls"
	goo_log "googo.io/goo-log"
	"io/ioutil"
)

type Tls struct {
	CaCrtFile     string
	ClientCrtFile string
	ClientKeyFile string
}

func (this *Tls) CaCrt() []byte {
	if this.CaCrtFile == "" {
		return caCert
	}
	bts, err := ioutil.ReadFile(this.CaCrtFile)
	if err != nil {
		goo_log.Error(err.Error())
	}
	return bts
}

func (this *Tls) ClientCrt() tls.Certificate {
	crt, err := tls.LoadX509KeyPair(this.ClientCrtFile, this.ClientKeyFile)
	if err != nil {
		goo_log.Error(err.Error())
	}
	return crt
}
