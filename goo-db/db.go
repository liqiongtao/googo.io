package goo_db

type DB interface {
	connect() (err error)
	ping()
}
