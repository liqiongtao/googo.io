package goo

import (
	goo_es "github.com/liqiongtao/googo.io/goo-es"
)

func ES() *goo_es.ESClient {
	return goo_es.Client()
}
