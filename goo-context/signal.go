package goo_context

import (
	"os"
	"os/signal"
	"syscall"
)

var (
	signalCH chan os.Signal
)

func init() {
	signalCH = make(chan os.Signal)
	signal.Notify(signalCH, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
}
