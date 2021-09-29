package goo

import (
	goo_db "googo.io/goo-db"
)

func DB() *goo_db.XOrm {
	return goo_db.XOrmClient()
}
