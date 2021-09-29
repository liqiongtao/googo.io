package goo_log

const (
	FileTag     = "tag"
	FileDirname = "dirname"
	FileMaxSize = "max-size"
)

func FileTagOption(tag string) Option {
	return NewOption(FileTag, tag)
}

func FileDirnameOption(dirname string) Option {
	return NewOption(FileDirname, dirname)
}

func FileMaxSizeOption(maxSize int64) Option {
	return NewOption(FileMaxSize, maxSize)
}
