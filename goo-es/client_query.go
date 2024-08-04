package goo_es

import (
	"bytes"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
	"io"
)

func (c *ESClient) Loop(index []string, filter []goo_utils.M, fn func(p goo_utils.Params) error) {
	size := 500
	for n := 0; n < 100; n++ {
		_, list := c.Query(index, filter, n*size, size)
		l := len(list)

		goo_log.Debug("[ES]", "查询数量", l, n, n*size, size)

		if l == 0 {
			return
		}

		for _, i := range list {
			if err := fn(i); err != nil {
				goo_log.Error("[ES]", err, i)
				return
			}
		}
	}
}

func (c *ESClient) Query(index []string, filter []goo_utils.M, offset, size int) (int64, []goo_utils.Params) {
	m := goo_utils.M{
		"from": offset,
		"size": size,
		"query": goo_utils.M{
			"bool": goo_utils.M{
				"filter": filter,
			},
		},
	}

	res, err := c.Search(index, m.Json())
	if err != nil {
		goo_log.Error("[ES]", err)
		return 0, []goo_utils.Params{}
	}
	defer res.Body.Close()

	var b bytes.Buffer
	if _, er := io.Copy(&b, res.Body); er != nil {
		goo_log.Error("[ES]", er.Error())
		return 0, []goo_utils.Params{}
	}

	p, er := goo_utils.Byte(b.Bytes()).Params()
	if er != nil {
		goo_log.Error("[ES]", er.Error())
		return 0, []goo_utils.Params{}
	}

	return p.Get("hits.total.value").Int64(), p.Get("hits.hits").Array()
}
