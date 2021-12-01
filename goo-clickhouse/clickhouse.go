package goo_clickhouse

import (
	"database/sql"
	"time"
)

var __client *client

func Init(config Config) {
	var err error
	if __client, err = New(config); err != nil {
		time.Sleep(10 * time.Second)
		Init(config)
	}
}

func DB() *sql.DB {
	return __client.db
}
