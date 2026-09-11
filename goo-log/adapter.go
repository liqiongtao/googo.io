package goolog

type Adapter interface {
	Write(msg *Message)
}
