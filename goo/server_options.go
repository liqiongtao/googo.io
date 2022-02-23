package goo

type Option struct {
	Name  string
	Value interface{}
}

func NewOption(name string, value interface{}) Option {
	return Option{Name: name, Value: value}
}

const (
	corsHeaderKeys = "cors-header-keys"
	noAccessPaths  = "no-access-paths"
	noLogPaths     = "no-log-paths"
)

func CorsHeaderKeysOption(headerKeys ...string) Option {
	return NewOption(corsHeaderKeys, headerKeys)
}

func NoAccessPathsOption(paths ...string) Option {
	return NewOption(noAccessPaths, paths)
}

func NoLogPathsOption(paths ...string) Option {
	return NewOption(noLogPaths, paths)
}
