package goo_request

import (
	"crypto/tls"
	"fmt"
	"os"
)

type Tls struct {
	CaCrtFile     string
	ClientCrtFile string
	ClientKeyFile string
}

func (s *Tls) CaCrt() ([]byte, error) {
	if s.CaCrtFile == "" {
		return caCert, nil
	}
	bts, err := os.ReadFile(s.CaCrtFile)
	if err != nil {
		return nil, fmt.Errorf("read ca cert: %w", err)
	}
	return bts, nil
}

func (s *Tls) ClientCrt() (tls.Certificate, error) {
	// 允许只配 CA、不做 mTLS
	if s.ClientCrtFile == "" || s.ClientKeyFile == "" {
		return tls.Certificate{}, nil
	}
	crt, err := tls.LoadX509KeyPair(s.ClientCrtFile, s.ClientKeyFile)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("load client cert: %w", err)
	}
	return crt, nil
}
