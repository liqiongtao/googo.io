package goo_log

type Option struct {
	Name  string
	Value interface{}
}

func NewOption(name string, value interface{}) Option {
	return Option{Name: name, Value: value}
}
