package goo_pprof

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	goo_log "github.com/liqiongtao/googo.io/goo-log"
)

type PProf struct {
	flag bool

	cpuFile string
	cpuFH   *os.File

	memoryFile string
	memoryFH   *os.File

	goroutineFile string
	goroutineFH   *os.File

	mutexFile string
	mutexFH   *os.File

	blockFile string
	blockFH   *os.File
}

func New(baseDir string) *PProf {
	if baseDir == "" {
		baseDir = "logs"
	}

	_ = os.MkdirAll(baseDir, 0755)

	ts := time.Now().Format("20060102150304")

	return &PProf{
		cpuFile:       fmt.Sprintf("%s/cpu.%s.prof", baseDir, ts),
		memoryFile:    fmt.Sprintf("%s/memory.%s.prof", baseDir, ts),
		goroutineFile: fmt.Sprintf("%s/goroutine.%s.prof", baseDir, ts),
		mutexFile:     fmt.Sprintf("%s/mutex.%s.prof", baseDir, ts),
		blockFile:     fmt.Sprintf("%s/block.%s.prof", baseDir, ts),
	}
}

func (pp *PProf) Start() {
	if pp.flag {
		goo_log.WithTag("goo-pprof").Info("正在执行")
		return
	}

	goo_log.WithTag("goo-pprof").Info("开始执行")

	pp.flag = true

	// 限制 CPU 使用数，避免过载
	runtime.GOMAXPROCS(1)
	// 开启对锁调用的跟踪
	runtime.SetMutexProfileFraction(1)
	// 开启对阻塞操作的跟踪
	runtime.SetBlockProfileRate(1)

	go pp.cpu()
	go pp.memory()
	go pp.goroutine()
	go pp.mutex()
	go pp.block()
}

func (pp *PProf) Stop() {
	pp.flag = false

	pprof.StopCPUProfile()

	if pp.cpuFH != nil {
		pp.cpuFH.Close()
	}

	if pp.memoryFH != nil {
		pp.memoryFH.Close()
	}

	if pp.goroutineFH != nil {
		pp.goroutineFH.Close()
	}

	if pp.mutexFH != nil {
		pp.mutexFH.Close()
	}

	if pp.blockFH != nil {
		pp.blockFH.Close()
	}

	time.Sleep(time.Second)

	goo_log.WithTag("goo-pprof").InfoF(
		"执行结束:\n%s\n%s\n%s\n%s\n%s",
		pp.memoryFile, pp.cpuFile, pp.goroutineFile, pp.blockFile, pp.mutexFile,
	)
}

func (pp *PProf) cpu() {
	var err error

	if pp.cpuFH, err = os.Create(pp.cpuFile); err != nil {
		goo_log.WithTag("goo-pprof").Error(err)
		return
	}

	pprof.StartCPUProfile(pp.cpuFH)
}

func (pp *PProf) memory() {
	var err error

	if pp.memoryFH, err = os.Create(pp.memoryFile); err != nil {
		goo_log.WithTag("goo-pprof").Error(err)
		return
	}

	runtime.GC()
	pprof.WriteHeapProfile(pp.memoryFH)
}

func (pp *PProf) goroutine() {
	var err error

	if pp.goroutineFH, err = os.Create(pp.goroutineFile); err != nil {
		goo_log.WithTag("goo-pprof").Error(err)
		return
	}

	if prof := pprof.Lookup("goroutine"); prof != nil {
		prof.WriteTo(pp.goroutineFH, 1)
	}
}

func (pp *PProf) mutex() {
	var err error

	if pp.mutexFH, err = os.Create(pp.mutexFile); err != nil {
		goo_log.WithTag("goo-pprof").Error(err)
		return
	}

	if prof := pprof.Lookup("mutex"); prof != nil {
		prof.WriteTo(pp.mutexFH, 1)
	}
}

func (pp *PProf) block() {
	var err error

	if pp.blockFH, err = os.Create(pp.blockFile); err != nil {
		goo_log.WithTag("goo-pprof").Error(err)
		return
	}

	if prof := pprof.Lookup("block"); prof != nil {
		prof.WriteTo(pp.blockFH, 1)
	}
}
