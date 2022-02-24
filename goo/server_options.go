package goo

type Option struct {
	Name  string
	Value interface{}
}

func NewOption(name string, value interface{}) Option {
	return Option{Name: name, Value: value}
}

const (
	baseDir       = "base-dir"
	serverName    = "server-name"
	corsHeaders   = "cors-headers"
	noAccessPaths = "no-access-paths"
	noLogPaths    = "no-log-paths"
)

// 根路径
func BaseDirOption(v string) Option {
	return NewOption(baseDir, v)
}

// 服务名称
func ServerNameOption(v string) Option {
	return NewOption(serverName, v)
}

// 跨域字段
func CorsHeadersOption(v ...string) Option {
	return NewOption(corsHeaders, v)
}

// 不允许访问的路径
func NoAccessPathsOption(v ...string) Option {
	return NewOption(noAccessPaths, v)
}

// 不记录日志的路径
func NoLogPathsOption(v ...string) Option {
	return NewOption(noLogPaths, v)
}
