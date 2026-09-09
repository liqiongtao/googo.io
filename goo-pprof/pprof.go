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

	// 开启对锁调用的跟踪
	runtime.SetMutexProfileFraction(1)
	// 开启对阻塞操作的跟踪
	runtime.SetBlockProfileRate(1)

	// CPU profile 必须同步启动，避免 StopCPUProfile 与异步 StartCPUProfile 竞态
	if err := pp.startCPU(); err != nil {
		goo_log.WithTag("goo-pprof").Error(err)
	}
	go pp.memory()
	go pp.goroutine()
	go pp.mutex()
	go pp.block()
}

func (pp *PProf) Stop() {
	pp.flag = false

	pprof.StopCPUProfile()
	runtime.SetMutexProfileFraction(0)
	runtime.SetBlockProfileRate(0)

	// 给异步 memory/goroutine/mutex/block 写盘一点时间，再关文件
	time.Sleep(200 * time.Millisecond)

	if pp.cpuFH != nil {
		_ = pp.cpuFH.Close()
		pp.cpuFH = nil
	}

	if pp.memoryFH != nil {
		_ = pp.memoryFH.Close()
		pp.memoryFH = nil
	}

	if pp.goroutineFH != nil {
		_ = pp.goroutineFH.Close()
		pp.goroutineFH = nil
	}

	if pp.mutexFH != nil {
		_ = pp.mutexFH.Close()
		pp.mutexFH = nil
	}

	if pp.blockFH != nil {
		_ = pp.blockFH.Close()
		pp.blockFH = nil
	}

	time.Sleep(time.Second)

	goo_log.WithTag("goo-pprof").InfoF(
		"执行结束:\n%s\n%s\n%s\n%s\n%s",
		pp.memoryFile, pp.cpuFile, pp.goroutineFile, pp.blockFile, pp.mutexFile,
	)
}

func (pp *PProf) startCPU() error {
	var err error
	if pp.cpuFH, err = os.Create(pp.cpuFile); err != nil {
		return err
	}
	if err = pprof.StartCPUProfile(pp.cpuFH); err != nil {
		_ = pp.cpuFH.Close()
		pp.cpuFH = nil
		return err
	}
	return nil
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
