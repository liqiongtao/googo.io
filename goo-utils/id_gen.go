package goo_utils

import (
	"strconv"
	"strings"
	"sync"
)

type iGenId interface {
	GenId() int64
}

var (
	__genId    iGenId
	__genIdOne sync.Once
)

func GenIdInit(adapter iGenId) {
	__genId = adapter
}

func GenId() int64 {
	__genIdOne.Do(func() {
		if __genId == nil {
			__genId = &SnowFlakeId{machineId: 1}
		}
	})
	return __genId.GenId()
}

func GenIdStr() string {
	return strings.ToLower(MD5([]byte(strconv.FormatInt(GenId(), 10))))
}
