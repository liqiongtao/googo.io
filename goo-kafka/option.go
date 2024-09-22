package goo_kafka

const (
	FocusName = "focus"
)

type Option struct {
	Name  string
	Value any
}

// 是否强制
func FocusOption() Option {
	return Option{Name: FocusName, Value: true}
}
