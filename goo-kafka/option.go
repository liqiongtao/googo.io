package goo_kafka

const (
	FocusName = "focus"
	RedisName = "redis"
)

type Option struct {
	Name  string
	Value interface{}
}

// 是否强制
func FocusOption() Option {
	return Option{Name: FocusName, Value: true}
}
