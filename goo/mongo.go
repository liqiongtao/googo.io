package goo

import goo_mongo "github.com/liqiongtao/googo.io/goo-mongo"

func Mongo(names ...string) *goo_mongo.Client {
	return goo_mongo.GetClient(names...)
}
