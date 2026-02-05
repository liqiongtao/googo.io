package goo_utils

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"time"
)

func SSLCertCheck(domain string) (*CertInfo, error) {
	info := &CertInfo{
		Domain: domain,
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         info.Domain,
	}

	conn, err := tls.DialWithDialer(
		&net.Dialer{Timeout: time.Second * 10},
		"tcp",
		fmt.Sprintf("%s:443", info.Domain),
		tlsConfig,
	)
	if err != nil {
		return nil, err
	}
	defer func() { _ = conn.Close() }()

	certChain := conn.ConnectionState().PeerCertificates
	if len(certChain) == 0 {
		return nil, errors.New("没有找到证书")
	}

	cert := certChain[0]

	// 填充证书信息
	info.Issuer = cert.Issuer.String()
	info.Subject = cert.Subject.String()
	info.NotBefore = cert.NotBefore
	info.NotAfter = cert.NotAfter
	info.SerialNumber = cert.SerialNumber.String()
	info.DNSNames = cert.DNSNames

	// 计算剩余天数
	now := time.Now()
	daysRemaining := int(cert.NotAfter.Sub(now).Hours() / 24)
	info.DaysRemaining = daysRemaining

	// 检查证书状态
	info.IsValid = now.After(cert.NotBefore) && now.Before(cert.NotAfter)
	info.IsExpired = now.After(cert.NotAfter)
	info.WillExpireSoon = daysRemaining <= 14 && daysRemaining > 0

	return info, nil
}

type CertInfo struct {
	Domain         string    `json:"domain"`
	Issuer         string    `json:"issuer"`           // 证书颁发者（CA机构）
	Subject        string    `json:"subject"`          // 证书持有者（域名所有者）
	NotBefore      time.Time `json:"not_before"`       // 证书生效开始时间
	NotAfter       time.Time `json:"not_after"`        // 证书过期时间
	DaysRemaining  int       `json:"days_remaining"`   // 距离证书过期的剩余天数, 负值表示已过期多少天，0表示今天过期
	IsValid        bool      `json:"is_valid"`         // 证书当前是否有效, false=证书已过期或未生效
	IsExpired      bool      `json:"is_expired"`       // 证书是否已过期, false=证书未过期
	WillExpireSoon bool      `json:"will_expire_soon"` // 证书是否即将过期, 小于14天时
	SerialNumber   string    `json:"serial_number"`    // 证书序列号, 由CA颁发的唯一标识符
	DNSNames       []string  `json:"dns_names"`        // 证书支持的域名列表（SAN - Subject Alternative Names）
}

func (info *CertInfo) Json() []byte {
	b, _ := json.Marshal(info)
	return b
}

func (info *CertInfo) String() string {
	return string(info.Json())
}
