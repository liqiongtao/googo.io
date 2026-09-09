package goo_utils

import (
	"strconv"
	"sync"

	"github.com/google/uuid"
)

type iGenId interface {
	GenId() int64
}

var (
	__genId    iGenId
	__genIdMu  sync.RWMutex
	__genIdOne sync.Once
)

func GenIdInit(adapter iGenId) {
	if adapter == nil {
		return
	}
	__genIdMu.Lock()
	__genId = adapter
	__genIdMu.Unlock()
}

func GenId() int64 {
	__genIdOne.Do(func() {
		__genIdMu.Lock()
		if __genId == nil {
			__genId = &SnowFlakeId{machineId: 1}
		}
		__genIdMu.Unlock()
	})
	__genIdMu.RLock()
	adapter := __genId
	__genIdMu.RUnlock()
	return adapter.GenId()
}

func GenIdStr() string {
	return strconv.FormatInt(GenId(), 10)
}

func UUID() string {
	return uuid.New().String()
}
