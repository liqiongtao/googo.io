package goo_utils

import (
	goo_log "github.com/liqiongtao/googo.io/goo-log"
)

func Recover() {
	if err := recover(); err != nil {
		goo_log.Error(err)
	}
}
