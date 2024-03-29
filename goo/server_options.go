package goo

var defaultOptions = &options{
	noAccessPath: map[string]struct{}{
		"/favicon.ico": {},
	},
	noLogPath: map[string]struct{}{
		"/favicon.ico": {},
	},
	corsHeaders: []string{
		"Content-Type", "Content-Length",
		"Accept", "Referer", "User-Agent", "Authorization",
		"X-Requested-Id", "X-Request-Timestamp", "X-Request-Sign",
		"X-Request-AppId", "X-Request-Source", "X-Request-Token",
	},
	disableEncryptionUris: map[string]struct{}{},
}

type options struct {
	pprofEnable bool

	serverName string
	env        Env

	corsHeaders  []string
	noAccessPath map[string]struct{}
	noLogPath    map[string]struct{}

	enableEncryption      bool
	encryptionKey         string
	disableEncryptionUris map[string]struct{}
}

type Option interface {
	apply(opts *options)
}

type funcOption struct {
	f func(opts *options)
}

func newFuncOption(f func(opts *options)) *funcOption {
	return &funcOption{f: f}
}

func (f funcOption) apply(opts *options) {
	f.f(opts)
}

// 开启分析
func PProfEnableOption(pprofEnable bool) Option {
	return newFuncOption(func(opts *options) {
		opts.pprofEnable = pprofEnable
	})
}

// 服务名称
func ServerNameOption(serverName string) Option {
	return newFuncOption(func(opts *options) {
		opts.serverName = serverName
	})
}

// 运行环境
func EnvOption(env Env) Option {
	return newFuncOption(func(opts *options) {
		opts.env = env
	})
}

// 跨域
func CorsHeaderOption(corsHeaders ...string) Option {
	return newFuncOption(func(opts *options) {
		opts.corsHeaders = append(opts.corsHeaders, corsHeaders...)
	})
}

// 禁止访问的path
func NoAccessPathsOption(noAccessPaths ...string) Option {
	return newFuncOption(func(opts *options) {
		for _, i := range noAccessPaths {
			opts.noAccessPath[i] = struct{}{}
		}
	})
}

// 不记录日志的path
func NoLogPathsOption(noLogPaths ...string) Option {
	return newFuncOption(func(opts *options) {
		for _, i := range noLogPaths {
			opts.noLogPath[i] = struct{}{}
		}
	})
}

// 启用加密传输
func EnableEncryptionOption(key ...string) Option {
	return newFuncOption(func(opts *options) {
		opts.enableEncryption = true
		if l := len(key); l > 0 {
			opts.encryptionKey = key[0]
		}
	})
}

// 启用加密传输
func EncryptionKeyOption(key string) Option {
	return newFuncOption(func(opts *options) {
		opts.encryptionKey = key
	})
}

// 不启用加密的urls
func DisableEncryptionUriOption(urls ...string) Option {
	return newFuncOption(func(opts *options) {
		for _, i := range urls {
			opts.disableEncryptionUris[i] = struct{}{}
		}
	})
}
