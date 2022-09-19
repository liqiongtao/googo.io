package goo_es

import (
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"strings"
)

type logger struct {
	level goo_log.Level
}

func (l logger) Printf(format string, v ...interface{}) {
	log := goo_log.WithTag("goo-es")
	switch l.level {
	case goo_log.DEBUG:
		log.DebugF(format, v...)
	case goo_log.INFO:
		log.InfoF(format, v...)
	case goo_log.WARN, goo_log.ERROR, goo_log.PANIC, goo_log.FATAL:
		if strings.Contains(format, "warning") {
			log.WarnF(format, v...)
			return
		}
		log.ErrorF(format, v...)
	}
}
