package goo_log

var defaultFileOptions = fileOptions{
	Filepath:     "logs/",
	MaxSize:      1 << 29,
	QueueSize:    0, // 0 表示按 CPU 自动计算
	DropWhenFull: false,
}

type fileOptions struct {
	Filepath     string
	MaxSize      int64
	QueueSize    int
	DropWhenFull bool
}

type FileOption interface {
	apply(*fileOptions)
}

type funcFileOption struct {
	f func(*fileOptions)
}

func (f *funcFileOption) apply(options *fileOptions) {
	f.f(options)
}

func newFuncFileOption(f func(*fileOptions)) *funcFileOption {
	return &funcFileOption{
		f: f,
	}
}

func FilePathOption(filepath string) FileOption {
	return newFuncFileOption(func(options *fileOptions) {
		options.Filepath = filepath
	})
}

func FileMaxSizeOption(maxSize int64) FileOption {
	return newFuncFileOption(func(options *fileOptions) {
		options.MaxSize = maxSize
	})
}

// FileQueueSizeOption 设置异步写队列长度；<=0 时按 CPU 自动计算。
func FileQueueSizeOption(size int) FileOption {
	return newFuncFileOption(func(options *fileOptions) {
		options.QueueSize = size
	})
}

// FileDropWhenFullOption 队列满时丢弃新日志且不阻塞调用方；默认 false（阻塞等待）。
func FileDropWhenFullOption(drop bool) FileOption {
	return newFuncFileOption(func(options *fileOptions) {
		options.DropWhenFull = drop
	})
}

