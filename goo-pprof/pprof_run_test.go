package goo_pprof

import (
	"testing"
)

func TestRegisterSignal(t *testing.T) {
	// 只验证可重复调用且不阻塞；真正的信号处理依赖进程信号，不在单测里等 Root().Done()
	RegisterSignal()
	RegisterSignal()
	Run()
}
