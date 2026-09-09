package goo_mail

import (
	"errors"
	"sync"
)

var (
	__mail iMail
	__mu   sync.RWMutex
)

func Init(conf Config) {
	m := New(conf)
	__mu.Lock()
	__mail = m
	__mu.Unlock()
}

func Send(msg Message) error {
	__mu.RLock()
	m := __mail
	__mu.RUnlock()
	if m == nil {
		return errors.New("mail not initialized, call goo_mail.Init first")
	}
	return m.Send(msg)
}
