package goo_utils

import (
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"github.com/shirou/gopsutil/v3/mem"
)

// 内存占比
func MemoryUsedPercent() (float64, error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		goo_log.Error(err)
		return 0, err
	}
	return v.UsedPercent, nil
}
