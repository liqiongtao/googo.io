package goo_es

import "github.com/olivere/elastic/v7"

var __client *client

func Init(conf Config, options ...elastic.ClientOptionFunc) (err error) {
	__client, err = New(conf, options...)
	return
}

func Default() *client {
	return __client
}

func Client() *elastic.Client {
	return __client.cli
}
