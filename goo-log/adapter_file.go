package goo_log

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"
)

// 文件日志适配器
type FileAdapter struct {
	file     *os.File
	date     string
	tag      string // 日志文件名标识
	filename string // 日志文件名称
	dirname  string // 日志保存目录
	count    int    // 日志文件切割数量
	maxSize  int64  // 日志文件切割大小
	mu       sync.Mutex
}

func NewFileAdapter(opts ...Option) *FileAdapter {
	fa := new(FileAdapter)

	if fa.dirname == "" {
		fa.dirname = "logs/"
	}

	if fa.maxSize == 0 {
		fa.maxSize = (1 << 20) * 21 // 100M
	}

	fa.SetOptions(opts...)

	if _, err := os.Stat(fa.dirname); os.IsNotExist(err) {
		os.MkdirAll(fa.dirname, 0755)
	}

	files, _ := filepath.Glob(fa.dirname + time.Now().Format("20060102") + "_*.log")
	fa.count = len(files)

	return fa
}

func (fa *FileAdapter) SetOptions(opts ...Option) {
	if ll := len(opts); ll == 0 {
		return
	}
	for _, opt := range opts {
		switch opt.Name {
		case FileTag:
			fa.tag = opt.Value.(string)
		case FileDirname:
			fa.dirname = path.Clean(opt.Value.(string)) + "/"
			if _, err := os.Stat(fa.dirname); os.IsNotExist(err) {
				os.MkdirAll(fa.dirname, 0755)
			}
		case FileMaxSize:
			fa.maxSize = opt.Value.(int64)
		}
	}
}

func (fa *FileAdapter) Write(msg *Message) {
	if w := fa.writer(); w != nil {
		w.Write(append(msg.JSON(), '\n'))
	}
}

func (fa *FileAdapter) writer() io.Writer {
	fa.mu.Lock()
	defer fa.mu.Unlock()

	// 定义当前日期
	fa.date = time.Now().Format("20060102")

	var (
		err      error
		filename string
	)

	// 定义文件名
	filename = fmt.Sprintf("%s%s.log", fa.dirname, fa.date)
	if fa.tag != "" {
		filename = fmt.Sprintf("%s%s_%s.log", fa.dirname, fa.date, fa.tag)
	}
	if fa.filename == "" {
		fa.filename = filename
	} else if fa.filename != filename {
		fa.count = 0
		fa.file = nil
		fa.filename = filename
	}

	// 获取文件属性
	fi, err := os.Stat(fa.filename)

	// 如果文件不存在
	if os.IsNotExist(err) {
		fa.file, err = os.OpenFile(fa.filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Println(err.Error())
		}
		return fa.file
	}

	// 如果文件不需要切割
	if fi.Size() < fa.maxSize {
		if fa.file == nil {
			fa.file, err = os.OpenFile(fa.filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				log.Println(err.Error())
			}
		}
		return fa.file
	}

	// 切割文件
	fa.cutFile()

	return fa.file
}

// 切割文件
func (fa *FileAdapter) cutFile() {
	fa.count += 1
	newFilename := fmt.Sprintf("%s%s_%d.log", fa.dirname, fa.date, fa.count)
	if fa.tag != "" {
		newFilename = fmt.Sprintf("%s%s_%s_%d.log", fa.dirname, fa.date, fa.tag, fa.count)
	}
	os.Rename(fa.filename, newFilename)
}
