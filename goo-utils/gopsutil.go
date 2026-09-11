package goo_utils

import (
	"errors"

	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

// MemoryUsedPercent 内存占比
func MemoryUsedPercent() (float64, error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		goo_log.Error(err)
		return 0, err
	}
	if v.Total == 0 {
		return 0, errors.New("memory total is 0")
	}
	return 100 - float64(v.Available)/float64(v.Total)*100, nil
}

// Memory 总内存 已用内存 可用内存 使用率
func Memory() (float64, float64, float64, float64, error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		goo_log.Error(err)
		return 0, 0, 0, 0, err
	}
	if v.Total == 0 {
		return 0, 0, 0, 0, errors.New("memory total is 0")
	}
	const gb = 1024 * 1024 * 1024
	return float64(v.Total) / gb,
		float64(v.Total-v.Available) / gb,
		float64(v.Available) / gb,
		100 - float64(v.Available)/float64(v.Total)*100,
		nil
}

// Disk 总大小 已用大小 可用大小 使用率
func Disk() (float64, float64, float64, float64, error) {
	v, err := disk.Usage("/")
	if err != nil {
		goo_log.Error(err)
		return 0, 0, 0, 0, err
	}
	const gb = 1024 * 1024 * 1024
	return float64(v.Total) / gb,
		float64(v.Used) / gb,
		float64(v.Free) / gb,
		v.UsedPercent,
		nil
}
