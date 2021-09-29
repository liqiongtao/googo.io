package goo_utils

import (
	goo_log "googo.io/goo-log"
)

func Recover() {
	if err := recover(); err != nil {
		goo_log.Error(err)
	}
}
