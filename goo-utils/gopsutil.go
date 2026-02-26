package goo_utils

import (
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"github.com/shirou/gopsutil/v3/disk"
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

// 内存
func Memory() (float64, float64, float64, error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		goo_log.Error(err)
		return 0, 0, 0, err
	}
	return float64(v.Total) / 1024 / 1024 / 1024,
		float64(v.Used) / 1024 / 1024 / 1024,
		float64(v.Available) / 1024 / 1024 / 1024,
		nil
}

// 空间
func Disk() (float64, float64, float64, error) {
	v, err := disk.Usage("/")
	if err != nil {
		goo_log.Error(err)
		return 0, 0, 0, err
	}
	return float64(v.Total) / 1024 / 1024 / 1024,
		float64(v.Used) / 1024 / 1024 / 1024,
		float64(v.Free) / 1024 / 1024 / 1024,
		nil
}
