package goocontext

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"

	goolog "github.com/liqiongtao/googo.io/goo-log"
)

// 信号分档（全进程只 Notify 一次；钩子异步派发，不堵接收循环）：
//
//	Exit:    SIGTERM / SIGINT / SIGQUIT(kill -3) → 先 cancel(Root)，再跑钩子；Wait() 等钩子结束
//	Restart: SIGHUP(kill -1)                     → 仅执行钩子，不 cancel
//	Action:  SIGUSR1 / SIGUSR2 等                → 仅执行钩子（如 pprof）

var (
	hookMu       sync.Mutex
	sigHooks     = make(map[os.Signal][]func())
	rootOnce     sync.Once
	exitOnce     sync.Once
	exiting      atomic.Bool
	rootCtx      context.Context
	rootCancel   context.CancelFunc
	exitFinished chan struct{} // Root() 后非 nil；退出钩子跑完后关闭
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

// OnExit 注册退出钩子（SIGTERM / SIGINT / SIGQUIT）。
// 先 cancel(Root) 再执行钩子，故钩子内可安全 <-Root().Done()。
// 主流程等关服完成用 Wait()；钩子内不要调用 Wait()（会等自己结束 → 死锁），用 Done() 或业务 WaitGroup。
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
					goolog.Error(r)
				}
			}()
			fn()
		}(fn)
	}
}

func isExitSignal(sig os.Signal) bool {
	return sig == syscall.SIGTERM || sig == syscall.SIGINT || sig == syscall.SIGQUIT
}

// Root 返回进程级共享 Context：全进程只监听一次信号。
// Done() 在收到退出信号时立刻触发；若要等 OnExit 钩子（如 Shutdown）结束，用 Wait()。
func Root() context.Context {
	rootOnce.Do(func() {
		rootCtx, rootCancel = context.WithCancel(context.Background())
		exitFinished = make(chan struct{})

		// 缓冲略大：钩子异步派发，避免短时连发信号被丢弃
		sig := make(chan os.Signal, 8)
		signal.Notify(sig,
			syscall.SIGHUP, syscall.SIGUSR1, syscall.SIGUSR2,
			syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT,
		)

		go func() {
			defer func() {
				if r := recover(); r != nil {
					goolog.Error(r)
				}
			}()

			for ch := range sig {
				s := ch
				if isExitSignal(s) {
					exitOnce.Do(func() {
						exiting.Store(true)
						go func() {
							defer close(exitFinished)
							rootCancel()
							runHooks(s)
						}()
					})
					continue
				}
				// 已进入退出流程则忽略 HUP/USR*，避免 Shutdown 中再 Upgrade
				if exiting.Load() {
					continue
				}
				go runHooks(s)
			}
		}()
	})
	return rootCtx
}

// Wait 阻塞直到退出钩子执行完毕（若尚未退出则一直等）。
// HTTP/gRPC 等主循环应 Wait()，以便 Shutdown/GracefulStop 完成后再返回。
// 勿在 OnExit 钩子内调用（会死锁）；钩子内请用 <-Root().Done()。
func Wait() {
	Root()
	<-exitFinished
}
