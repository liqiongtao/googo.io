package goo_pprof

import (
	"os"
	"sync"
	"syscall"

	"github.com/liqiongtao/googo.io/goo-context"
	goolog "github.com/liqiongtao/googo.io/goo-log"
)

var (
	ppMu       sync.Mutex
	currentPP  *PProf
	signalOnce sync.Once
)

// StartDefault 启动进程级默认 pprof（若已在跑则忽略）。
func StartDefault() {
	ppMu.Lock()
	defer ppMu.Unlock()
	if currentPP != nil {
		return
	}
	pp := New("logs")
	if err := pp.Start(); err != nil {
		goolog.WithTag("goo-pprof").Error(err)
		return
	}
	currentPP = pp
}

// StopDefault 停止进程级默认 pprof。
func StopDefault() {
	ppMu.Lock()
	defer ppMu.Unlock()
	if currentPP == nil {
		return
	}
	currentPP.Stop()
	currentPP = nil
}

// Toggle 切换进程级默认 pprof 开/关（持锁串行，避免连发 USR1 交叉 Start/Stop）。
func Toggle() {
	ppMu.Lock()
	defer ppMu.Unlock()
	if currentPP != nil {
		currentPP.Stop()
		currentPP = nil
		goolog.WithTag("goo-pprof").Info("pprof 已停止")
		return
	}
	pp := New("logs")
	if err := pp.Start(); err != nil {
		goolog.WithTag("goo-pprof").Error(err)
		return
	}
	currentPP = pp
	goolog.WithTag("goo-pprof").Info("pprof 已开始")
}

// RegisterSignal 全进程只注册一次：SIGUSR1 → Toggle。
func RegisterSignal() {
	signalOnce.Do(func() {
		goocontext.OnSignal(syscall.SIGUSR1, Toggle)
		goolog.InfoF("pprof 已注册，切换分析: kill -USR1 %d", os.Getpid())
	})
}

// Run 等价于 RegisterSignal（保留旧名）。
func Run() {
	RegisterSignal()
}
