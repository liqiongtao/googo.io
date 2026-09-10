package goo_es

import (
	"bytes"
	"io"

	goo_log "github.com/liqiongtao/googo.io/goo-log"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
)

func (c *ESClient) Loop(index []string, filter []goo_utils.M, fn func(p goo_utils.Params) error) {
	size := 500
	for n := 0; n < 1000; n++ {
		_, list := c.Query(index, filter, n*size, size)
		l := len(list)

		if c.showLog {
			goo_log.Debug("[ES]", "查询数量", l, n, n*size, size)
		}

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

	if res.IsError() {
		goo_log.Error("[ES]", res.String())
		return 0, []goo_utils.Params{}
	}

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

func (c *ESClient) LoopV2(index []string, m goo_utils.M, fn func(p goo_utils.Params) error) {
	body := goo_utils.M{}
	for k, v := range m {
		body[k] = v
	}
	size := 500
	for n := 0; n < 1000; n++ {
		body["from"] = n * size
		body["size"] = size
		_, list := c.QueryV2(index, body)
		l := len(list)

		if c.showLog {
			goo_log.Debug("[ES]", "查询数量", l, n, n*size, size)
		}

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

func (c *ESClient) QueryV2(index []string, m goo_utils.M) (int64, []goo_utils.Params) {
	res, err := c.Search(index, m.Json())
	if err != nil {
		goo_log.Error("[ES]", err)
		return 0, []goo_utils.Params{}
	}
	defer res.Body.Close()

	if res.IsError() {
		goo_log.Error("[ES]", res.String())
		return 0, []goo_utils.Params{}
	}

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
