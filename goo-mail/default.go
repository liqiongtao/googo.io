package goo_mail

import "errors"

var (
	__mail iMail
)

func Init(conf Config) {
	__mail = New(conf)
}

func Send(msg Message) error {
	if __mail == nil {
		return errors.New("mail not initialized, call goo_mail.Init first")
	}
	return __mail.Send(msg)
}
