package gooes

import (
	"bytes"
	"io"

	goolog "github.com/liqiongtao/googo.io/goo-log"
	gooutils "github.com/liqiongtao/googo.io/goo-util"
)

func (c *ESClient) Loop(index []string, filter []gooutils.M, fn func(p gooutils.Params) error) {
	size := 500
	for n := 0; n < 1000; n++ {
		_, list := c.Query(index, filter, n*size, size)
		l := len(list)

		if c.showLog {
			goolog.Debug("[ES]", "查询数量", l, n, n*size, size)
		}

		if l == 0 {
			return
		}

		for _, i := range list {
			if err := fn(i); err != nil {
				goolog.Error("[ES]", err, i)
				return
			}
		}
	}
}

func (c *ESClient) Query(index []string, filter []gooutils.M, offset, size int) (int64, []gooutils.Params) {
	m := gooutils.M{
		"from": offset,
		"size": size,
		"query": gooutils.M{
			"bool": gooutils.M{
				"filter": filter,
			},
		},
	}

	res, err := c.Search(index, m.Json())
	if err != nil {
		goolog.Error("[ES]", err)
		return 0, []gooutils.Params{}
	}
	defer res.Body.Close()

	if res.IsError() {
		goolog.Error("[ES]", res.String())
		return 0, []gooutils.Params{}
	}

	var b bytes.Buffer
	if _, er := io.Copy(&b, res.Body); er != nil {
		goolog.Error("[ES]", er.Error())
		return 0, []gooutils.Params{}
	}

	p, er := gooutils.Byte(b.Bytes()).Params()
	if er != nil {
		goolog.Error("[ES]", er.Error())
		return 0, []gooutils.Params{}
	}

	return p.Get("hits.total.value").Int64(), p.Get("hits.hits").Array()
}

func (c *ESClient) LoopV2(index []string, m gooutils.M, fn func(p gooutils.Params) error) {
	body := gooutils.M{}
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
			goolog.Debug("[ES]", "查询数量", l, n, n*size, size)
		}

		if l == 0 {
			return
		}

		for _, i := range list {
			if err := fn(i); err != nil {
				goolog.Error("[ES]", err, i)
				return
			}
		}
	}
}

func (c *ESClient) QueryV2(index []string, m gooutils.M) (int64, []gooutils.Params) {
	res, err := c.Search(index, m.Json())
	if err != nil {
		goolog.Error("[ES]", err)
		return 0, []gooutils.Params{}
	}
	defer res.Body.Close()

	if res.IsError() {
		goolog.Error("[ES]", res.String())
		return 0, []gooutils.Params{}
	}

	var b bytes.Buffer
	if _, er := io.Copy(&b, res.Body); er != nil {
		goolog.Error("[ES]", er.Error())
		return 0, []gooutils.Params{}
	}

	p, er := gooutils.Byte(b.Bytes()).Params()
	if er != nil {
		goolog.Error("[ES]", er.Error())
		return 0, []gooutils.Params{}
	}

	return p.Get("hits.total.value").Int64(), p.Get("hits.hits").Array()
}
