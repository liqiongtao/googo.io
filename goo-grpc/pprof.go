package goo_grpc

import (
	"fmt"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
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

func newPProf() *PProf {
	ts := time.Now().Unix()
	return &PProf{
		cpuFile:       fmt.Sprintf("cpu.%d.prof", ts),
		memoryFile:    fmt.Sprintf("memory.%d.prof", ts),
		goroutineFile: fmt.Sprintf("goroutine.%d.prof", ts),
		mutexFile:     fmt.Sprintf("mutex.%d.prof", ts),
		blockFile:     fmt.Sprintf("block.%d.prof", ts),
	}
}

func (pp *PProf) start() {
	if pp.flag {
		return
	}

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

func (pp *PProf) cpu() {
	var err error

	if pp.cpuFH, err = os.Create(pp.cpuFile); err != nil {
		goo_log.WithTag("goo-grpc").Error(err)
		return
	}

	pprof.StartCPUProfile(pp.cpuFH)
}

func (pp *PProf) memory() {
	var err error

	if pp.memoryFH, err = os.Create(pp.memoryFile); err != nil {
		goo_log.WithTag("goo-grpc").Error(err)
		return
	}

	runtime.GC()
	pprof.WriteHeapProfile(pp.memoryFH)
}

func (pp *PProf) goroutine() {
	var err error

	if pp.goroutineFH, err = os.Create(pp.goroutineFile); err != nil {
		goo_log.WithTag("goo-grpc").Error(err)
		return
	}

	if prof := pprof.Lookup("goroutine"); prof != nil {
		prof.WriteTo(pp.goroutineFH, 1)
	}
}

func (pp *PProf) mutex() {
	var err error

	if pp.mutexFH, err = os.Create(pp.mutexFile); err != nil {
		goo_log.WithTag("goo-grpc").Error(err)
		return
	}

	if prof := pprof.Lookup("mutex"); prof != nil {
		prof.WriteTo(pp.mutexFH, 1)
	}
}

func (pp *PProf) block() {
	var err error

	if pp.blockFH, err = os.Create(pp.blockFile); err != nil {
		goo_log.WithTag("goo-grpc").Error(err)
		return
	}

	if prof := pprof.Lookup("block"); prof != nil {
		prof.WriteTo(pp.blockFH, 1)
	}
}

func (pp *PProf) stop() {
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
}
