package goocontext

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	goo_log "github.com/liqiongtao/googo.io/goo-log"
)

// 信号分档（全进程只 Notify 一次；钩子异步派发，不堵接收循环）：
//
//	Exit:    SIGTERM / SIGINT / SIGQUIT(kill -3) → 钩子后 cancel(Root)，只一次
//	Restart: SIGHUP(kill -1)                     → 仅执行钩子，不 cancel
//	Action:  SIGUSR1 / SIGUSR2 等                → 仅执行钩子（如 pprof）

var (
	hookMu     sync.Mutex
	sigHooks   = make(map[os.Signal][]func())
	rootOnce   sync.Once
	exitOnce   sync.Once
	exiting    atomic.Bool
	rootCtx    context.Context
	rootCancel context.CancelFunc
)

// OnSignal 注册某信号的自定义逻辑（进程级，不提供注销）。
func OnSignal(sig os.Signal, fn func()) {
	if fn == nil {
		return
	}
	hookMu.Lock()
	sigHooks[sig] = append(sigHooks[sig], fn)
	hookMu.Unlock()
	Root()
}

// OnExit 注册退出钩子（SIGTERM / SIGINT / SIGQUIT），在 cancel(Root) 之前执行。
func OnExit(fn func()) {
	OnSignal(syscall.SIGTERM, fn) // kill -15
	OnSignal(syscall.SIGINT, fn)  // kill -2
	OnSignal(syscall.SIGQUIT, fn) // kill -3
}

// OnRestart 注册重启钩子（SIGHUP / kill -1），不取消 Root。
func OnRestart(fn func()) {
	OnSignal(syscall.SIGHUP, fn)
}

func runHooks(sig os.Signal) {
	hookMu.Lock()
	fns := make([]func(), len(sigHooks[sig]))
	copy(fns, sigHooks[sig])
	hookMu.Unlock()

	for _, fn := range fns {
		func(fn func()) {
			defer func() {
				if r := recover(); r != nil {
					goo_log.Error(r)
				}
			}()
			fn()
		}(fn)
	}
}

func isExitSignal(sig os.Signal) bool {
	return sig == syscall.SIGTERM || sig == syscall.SIGINT || sig == syscall.SIGQUIT
}

// NotifyParentExit 子进程继承 listener 启动成功后，通知父进程退出（grace handoff）。
func NotifyParentExit() {
	if os.Getenv("LISTEN_FDS") == "" {
		return
	}
	ppid := os.Getppid()
	if ppid <= 1 {
		return
	}
	if err := syscall.Kill(ppid, syscall.SIGTERM); err != nil {
		goo_log.ErrorF("NotifyParentExit kill ppid=%d err=%v", ppid, err)
	}
}

// NotifyParentExitAfter 延迟通知父进程退出，给 Serve 进入 Accept 留出时间。
func NotifyParentExitAfter(d time.Duration) {
	if os.Getenv("LISTEN_FDS") == "" {
		return
	}
	time.AfterFunc(d, NotifyParentExit)
}

// Root 返回进程级共享 Context：全进程只监听一次信号。
func Root() context.Context {
	rootOnce.Do(func() {
		rootCtx, rootCancel = context.WithCancel(context.Background())

		// 缓冲略大：钩子异步派发，避免短时连发信号被丢弃
		sig := make(chan os.Signal, 8)
		signal.Notify(sig,
			syscall.SIGHUP, syscall.SIGUSR1, syscall.SIGUSR2,
			syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT,
		)

		go func() {
			defer func() {
				if r := recover(); r != nil {
					goo_log.Error(r)
				}
			}()

			for ch := range sig {
				s := ch
				if isExitSignal(s) {
					// 只退出一次；钩子在独立 goroutine，不堵信号循环
					// cancel 仍在钩子之后，保证 <-Root().Done() 等 Shutdown 完成
					exitOnce.Do(func() {
						exiting.Store(true)
						go func() {
							runHooks(s)
							rootCancel()
						}()
					})
					continue
				}
				// 已进入退出流程则忽略 HUP/USR*，避免 Shutdown 中再 StartProcess
				if exiting.Load() {
					continue
				}
				go runHooks(s)
			}
		}()
	})
	return rootCtx
}
