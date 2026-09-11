package goo

import (
	goodb "github.com/liqiongtao/googo.io/goo-db"
)

func DB(names ...string) *goodb.Client {
	return goodb.GetClient(names...)
}
