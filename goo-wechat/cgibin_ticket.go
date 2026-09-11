package goowechat

import (
	"encoding/json"
	"errors"
	"fmt"
	goolog "github.com/liqiongtao/googo.io/goo-log"
	goorequest "github.com/liqiongtao/googo.io/goo-request"
	"time"
)

type cgiTicket struct {
	Appid  string
	Secret string
}

func CGITicket(appid, secret string) *cgiTicket {
	return &cgiTicket{Appid: appid, Secret: secret}
}

func (this *cgiTicket) Get() string {
	key := fmt.Sprintf(cgi_ticket_key, this.Appid)
	return __cache.Get(key).Val()
}

func (this *cgiTicket) TTL() time.Duration {
	key := fmt.Sprintf(cgi_ticket_key, this.Appid)
	return __cache.TTL(key).Val()
}

func (this *cgiTicket) Set() error {
	accessToken := CGIToken(this.Appid, this.Secret).Get()
	buf, _ := goorequest.Get(fmt.Sprintf(cgi_ticket_url, accessToken))

	rsp := struct {
		Ticket    string `json:"ticket"`
		ExpiresIn int64  `json:"expires_in"`
		ErrCode   int    `json:"errcode"`
		ErrMsg    string `json:"errmsg"`
	}{}

	if err := json.Unmarshal(buf, &rsp); err != nil {
		return err
	}
	if errCode := rsp.ErrCode; errCode != 0 {
		return errors.New(rsp.ErrMsg)
	}

	goolog.WithField("ticket", rsp.Ticket).WithField("expire_in", rsp.ExpiresIn).Debug()

	key := fmt.Sprintf(cgi_ticket_key, this.Appid)
	return __cache.Set(key, rsp.Ticket, time.Duration(rsp.ExpiresIn)*time.Second).Err()
}
