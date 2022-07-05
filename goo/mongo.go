package goo

import (
	goo_mongo "github.com/liqiongtao/googo.io/goo-mongo"
	"go.mongodb.org/mongo-driver/mongo"
)

func Mongo(names ...string) *mongo.Database {
	return goo_mongo.GetClient(names...).DB()
}

func MongoClient(names ...string) *goo_mongo.Client {
	return goo_mongo.GetClient(names...)
}
