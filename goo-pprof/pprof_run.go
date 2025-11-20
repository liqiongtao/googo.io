package goo_pprof

import (
	"os"
	"os/signal"
	"syscall"

	goo_log "github.com/liqiongtao/googo.io/goo-log"
)

func Run() {
	pp := New()

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGUSR1, syscall.SIGUSR2)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				goo_log.Error(r)
			}
		}()

		for ch := range sig {
			switch ch {
			case syscall.SIGUSR1: // kill -USR1
				pp.Start()

			case syscall.SIGUSR2: // kill -USR2
				pp.Stop()
			}
		}
	}()

	goo_log.InfoF("pprof 已启动:\n开始执行分析: kill -USR1 %d\n结束执行分析: kill -USR2 %d\n", os.Getpid(), os.Getpid())
}
