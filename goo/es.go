package goo

import (
	gooes "github.com/liqiongtao/googo.io/goo-es"
)

func ES() *gooes.ESClient {
	return gooes.Client()
}
