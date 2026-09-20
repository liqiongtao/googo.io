package goo

import (
	goomongo "github.com/liqiongtao/googo.io/goo-mongo"
	"go.mongodb.org/mongo-driver/mongo"
)

func Mongo(names ...string) *mongo.Database {
	cli := goomongo.GetClient(names...)
	if cli == nil {
		return nil
	}
	return cli.DB()
}

func MongoClient(names ...string) *goomongo.Client {
	return goomongo.GetClient(names...)
}
