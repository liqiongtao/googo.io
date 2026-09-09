package goo_utils

import (
	"sync"
	"time"
)

// 雪花算法
type SnowFlakeId struct {
	dataCenterId int // 机房ID
	machineId    int // 机器ID

	lastTime int64 // 最后时间
	sn       int   // 序号

	mu sync.Mutex
}

func (sf *SnowFlakeId) waitNextMillis(last int64) int64 {
	ts := time.Now().UnixNano() / 1e6
	for ts <= last {
		time.Sleep(time.Millisecond)
		ts = time.Now().UnixNano() / 1e6
	}
	return ts
}

func (sf *SnowFlakeId) GenId() int64 {
	sf.mu.Lock()
	defer sf.mu.Unlock()

	for {
		ts := time.Now().UnixNano() / 1e6

		// 时钟回拨：解锁等待，避免堵住所有发号
		if ts < sf.lastTime {
			last := sf.lastTime
			sf.mu.Unlock()
			_ = sf.waitNextMillis(last)
			sf.mu.Lock()
			continue
		}

		if sf.lastTime == ts {
			// 2^12 - 1 = 4095，每毫秒最多 4096 个 ID（0~4095）
			if sf.sn >= 4095 {
				last := sf.lastTime
				sf.mu.Unlock()
				_ = sf.waitNextMillis(last)
				sf.mu.Lock()
				continue
			}
			sf.sn++
		} else {
			sf.sn = 0
		}

		sf.lastTime = ts

		// 机房ID / 机器ID 做掩码，避免越界污染序号位
		dataCenterId := (sf.dataCenterId & 0x1f) << 17
		machineId := (sf.machineId & 0x1f) << 12

		return (ts << 22) | int64(dataCenterId) | int64(machineId) | int64(sf.sn)
	}
}
