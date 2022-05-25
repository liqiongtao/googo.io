package goo_log

func FilePathOption(filepath string) Option {
	return NewOption("filepath", filepath)
}

func FileMaxSizeOption(maxSize int64) Option {
	return NewOption("max-size", maxSize)
}
