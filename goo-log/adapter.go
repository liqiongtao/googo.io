package goo_log

type Adapter interface {
	Write(msg *Message)
}
