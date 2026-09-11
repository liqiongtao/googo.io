package goo_pprof

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"sync"
	"time"

	goolog "github.com/liqiongtao/googo.io/goo-log"
)

type PProf struct {
	mu   sync.Mutex
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

func (pp *PProf) Start() error {
	pp.mu.Lock()
	defer pp.mu.Unlock()

	if pp.flag {
		goolog.WithTag("goo-pprof").Info("正在执行")
		return nil
	}

	goolog.WithTag("goo-pprof").Info("开始执行")

	// 开启对锁调用的跟踪
	runtime.SetMutexProfileFraction(1)
	// 开启对阻塞操作的跟踪
	runtime.SetBlockProfileRate(1)

	// 仅启动 CPU 采样；memory/goroutine/mutex/block 在 Stop 时写入，此时才有有效样本
	if err := pp.startCPU(); err != nil {
		goolog.WithTag("goo-pprof").Error(err)
		runtime.SetMutexProfileFraction(0)
		runtime.SetBlockProfileRate(0)
		return err
	}

	pp.flag = true
	return nil
}

func (pp *PProf) Stop() {
	pp.mu.Lock()
	defer pp.mu.Unlock()

	if !pp.flag {
		return
	}
	pp.flag = false

	pprof.StopCPUProfile()

	// 在关闭 mutex/block 采样前落盘，保留 Start～Stop 期间的样本
	pp.writeMemory()
	pp.writeGoroutine()
	pp.writeMutex()
	pp.writeBlock()

	runtime.SetMutexProfileFraction(0)
	runtime.SetBlockProfileRate(0)

	pp.closeFH(&pp.cpuFH)
	pp.closeFH(&pp.memoryFH)
	pp.closeFH(&pp.goroutineFH)
	pp.closeFH(&pp.mutexFH)
	pp.closeFH(&pp.blockFH)

	goolog.WithTag("goo-pprof").InfoF(
		"执行结束:\n%s\n%s\n%s\n%s\n%s",
		pp.memoryFile, pp.cpuFile, pp.goroutineFile, pp.blockFile, pp.mutexFile,
	)
}

func (pp *PProf) closeFH(fh **os.File) {
	if fh == nil || *fh == nil {
		return
	}
	_ = (*fh).Close()
	*fh = nil
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

func (pp *PProf) writeMemory() {
	var err error
	if pp.memoryFH, err = os.Create(pp.memoryFile); err != nil {
		goolog.WithTag("goo-pprof").Error(err)
		return
	}
	runtime.GC()
	_ = pprof.WriteHeapProfile(pp.memoryFH)
}

func (pp *PProf) writeGoroutine() {
	var err error
	if pp.goroutineFH, err = os.Create(pp.goroutineFile); err != nil {
		goolog.WithTag("goo-pprof").Error(err)
		return
	}
	if prof := pprof.Lookup("goroutine"); prof != nil {
		_ = prof.WriteTo(pp.goroutineFH, 1)
	}
}

func (pp *PProf) writeMutex() {
	var err error
	if pp.mutexFH, err = os.Create(pp.mutexFile); err != nil {
		goolog.WithTag("goo-pprof").Error(err)
		return
	}
	if prof := pprof.Lookup("mutex"); prof != nil {
		_ = prof.WriteTo(pp.mutexFH, 1)
	}
}

func (pp *PProf) writeBlock() {
	var err error
	if pp.blockFH, err = os.Create(pp.blockFile); err != nil {
		goolog.WithTag("goo-pprof").Error(err)
		return
	}
	if prof := pprof.Lookup("block"); prof != nil {
		_ = prof.WriteTo(pp.blockFH, 1)
	}
}
