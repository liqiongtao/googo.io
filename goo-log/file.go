package goo_log

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type fileJob struct {
	data     []byte
	syncDone chan error
}

type FileAdapter struct {
	filepath string
	filename string
	ymd      string
	fh       *os.File

	maxSize      int64
	count        int
	dropWhenFull bool
	dropped      atomic.Uint64

	ch     chan fileJob
	mu     sync.Mutex
	wg     sync.WaitGroup
	closed atomic.Bool
}

func NewFileLog(opts ...FileOption) *Logger {
	return New(NewFileAdapter(opts...))
}

func NewFileAdapter(opt ...FileOption) *FileAdapter {
	opts := defaultFileOptions
	for _, o := range opt {
		o.apply(&opts)
	}

	queueSize := opts.QueueSize
	if queueSize <= 0 {
		queueSize = runtime.NumCPU() * 64
		if queueSize < 64 {
			queueSize = 64
		}
	}

	fa := &FileAdapter{
		filepath:     normalizeFilePath(opts.Filepath),
		maxSize:      opts.MaxSize,
		dropWhenFull: opts.DropWhenFull,
		ch:           make(chan fileJob, queueSize),
	}

	if err := os.MkdirAll(fa.filepath, 0755); err != nil {
		log.Println("[goo-log][file]", err.Error())
	}

	fa.ymd = time.Now().Format("20060102")
	fa.count = fa.rotateCount(fa.ymd)

	fa.wg.Add(1)
	go fa.loop()

	return fa
}

func normalizeFilePath(path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		path = defaultFileOptions.Filepath
	}
	return strings.TrimRight(path, "/") + "/"
}

func (fa *FileAdapter) loop() {
	defer fa.wg.Done()

	for job := range fa.ch {
		if job.syncDone != nil {
			job.syncDone <- fa.fsync()
			continue
		}
		fa.writeHandle(job.data)
	}

	_ = fa.fsync()
	fa.mu.Lock()
	fa.closeFile()
	fa.mu.Unlock()
}

func (fa *FileAdapter) fsync() error {
	fa.mu.Lock()
	defer fa.mu.Unlock()

	if fa.fh == nil {
		return nil
	}
	if err := fa.fh.Sync(); err != nil {
		log.Println("[goo-log][file]", err.Error())
		return err
	}
	return nil
}

func (fa *FileAdapter) Write(msg *Message) {
	if msg == nil || msg.Entry == nil || fa.closed.Load() {
		return
	}

	b := msg.JSON()
	if len(b) == 0 {
		return
	}

	job := fileJob{data: b}
	if fa.dropWhenFull {
		defer func() { recover() }()
		select {
		case fa.ch <- job:
		default:
			n := fa.dropped.Add(1)
			if n == 1 || n%1000 == 0 {
				log.Printf("[goo-log][file] queue full, dropped %d logs", n)
			}
		}
		return
	}

	// Close 之后再 send 会 panic，仅保护 channel send
	defer func() { recover() }()
	fa.ch <- job
}

// Dropped 返回因队列满而丢弃的日志条数（仅 DropWhenFull 模式有意义）。
func (fa *FileAdapter) Dropped() uint64 {
	return fa.dropped.Load()
}

func (fa *FileAdapter) Sync() error {
	if fa.closed.Load() {
		return nil
	}

	done := make(chan error, 1)
	sent := true
	func() {
		defer func() {
			if recover() != nil {
				sent = false
			}
		}()
		fa.ch <- fileJob{syncDone: done}
	}()
	if !sent {
		return nil
	}
	return <-done
}

func (fa *FileAdapter) Close() error {
	if !fa.closed.CompareAndSwap(false, true) {
		return nil
	}

	close(fa.ch)
	fa.wg.Wait()
	return nil
}

func (fa *FileAdapter) writeHandle(b []byte) {
	if len(b) == 0 {
		return
	}

	fa.mu.Lock()
	defer fa.mu.Unlock()

	var (
		nw       = time.Now()
		ymd      = nw.Format("20060102")
		filename = ymd + ".log"
	)

	if filename != fa.filename {
		if ymd != fa.ymd {
			fa.ymd = ymd
			fa.count = fa.rotateCount(ymd)
		}
		fa.filename = filename
		fa.closeFile()
	}

	if fa.fh == nil {
		if err := fa.openFile(); err != nil {
			return
		}
	}

	if fa.maxSize > 0 {
		// 切割失败时仍尽量写入当前句柄，避免丢日志
		_ = fa.cutFile(ymd)
	}

	if fa.fh == nil {
		if err := fa.openFile(); err != nil {
			return
		}
	}

	if _, err := fa.fh.Write(b); err != nil {
		log.Println("[goo-log][file]", err.Error())
		return
	}
	if _, err := fa.fh.Write([]byte("\n")); err != nil {
		log.Println("[goo-log][file]", err.Error())
	}
}

func (fa *FileAdapter) closeFile() {
	if fa.fh == nil {
		return
	}

	if err := fa.fh.Sync(); err != nil {
		log.Println("[goo-log][file]", err.Error())
	}
	if err := fa.fh.Close(); err != nil {
		log.Println("[goo-log][file]", err.Error())
	}
	fa.fh = nil
}

func (fa *FileAdapter) openFile() (err error) {
	fa.fh, err = os.OpenFile(fa.filepath+fa.filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("[goo-log][file]", err.Error())
	}
	return
}

func (fa *FileAdapter) cutFile(ymd string) (err error) {
	var info os.FileInfo

	info, err = os.Stat(fa.filepath + fa.filename)
	if err != nil {
		log.Println("[goo-log][file]", err.Error())
		return
	}
	if info.Size() < fa.maxSize {
		return
	}

	// 先找可用序号，成功 rename 后再提交 fa.count，避免失败时序号被抬高
	next := fa.count
	var filename string
	for {
		next++
		if next <= 0 || next-fa.count > 1000000 {
			err = fmt.Errorf("rotate sequence overflow: count=%d next=%d", fa.count, next)
			log.Println("[goo-log][file]", err.Error())
			return
		}
		filename = fmt.Sprintf("%s_%d.log", fa.filepath+ymd, next)
		_, err = os.Stat(filename)
		if os.IsNotExist(err) {
			break
		}
		if err != nil {
			log.Println("[goo-log][file]", err.Error())
			return
		}
	}

	if err = os.Rename(fa.filepath+fa.filename, filename); err != nil {
		log.Println("[goo-log][file]", err.Error())
		return
	}
	fa.count = next

	fa.closeFile()

	if err = fa.openFile(); err != nil {
		return
	}
	return
}

func (fa *FileAdapter) rotateCount(ymd string) int {
	files, err := filepath.Glob(fa.filepath + ymd + "_*.log")
	if err != nil {
		log.Println("[goo-log][file]", err.Error())
		return 0
	}

	// 取已有切割文件的最大序号，而不是文件个数（避免空洞时错位，或误判）
	max := 0
	for _, f := range files {
		base := strings.TrimSuffix(filepath.Base(f), ".log")
		i := strings.LastIndex(base, "_")
		if i < 0 || i+1 >= len(base) {
			continue
		}
		n, convErr := strconv.Atoi(base[i+1:])
		if convErr != nil {
			continue
		}
		if n > max {
			max = n
		}
	}
	return max
}
