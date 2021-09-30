package goo

import (
	goo_db "github.com/liqiongtao/googo.io/goo-db"
)

func DB(names ...string) *goo_db.XOrm {
	return goo_db.XOrmClient(names...)
}
