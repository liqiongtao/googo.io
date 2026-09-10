package goo_es

import (
	"bytes"
	"context"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
	"io"
	"time"
)

// 搜索
//
// - range: 范围搜索
// - match: 匹配搜索
// - match_phrase: 全文检索
// - nested: 嵌套搜索
// - inner_hits: 只返回匹配项
// - inner_hits.size: 只返回匹配项的第1项
// - highlight: 高亮
// - _source.excludes: 不返回字段
//
// GET /abd-*/_search
//
//	{
//	  "from" : 0,
//	  "size" : 10,
//	  "query": {
//	    "bool": {
//	      "must": [
//	        {"match": { "user_id": 2716 }},
//	        {"range": { "datetime": { "gte": "2024-01-14 00:00:00", "lte": "2024-01-14 23:59:59" }}},
//	        {"range": { "size": { "gte": 60, "lte": 3600 }}},
//	        {
//	          "nested": {
//	            "path": "list",
//	            "query": {
//	              "match_phrase": { "list.title": "青岛" }
//	            },
//	            "inner_hits": {
//	              "size": 1
//	            }
//	          }
//	        }
//	      ]
//	    }
//	  },
//	  "highlight": {
//	    "pre_tags": "<em>",
//	    "post_tags": "</em>",
//	    "fields": {
//	      "list.title": {}
//	    }
//	  },
//	  "_source": {
//	    "excludes" : ["list"]
//	  }
//	}
func (c *ESClient) Search(index []string, body []byte) (*esapi.Response, error) {
	req := esapi.SearchRequest{
		Index: index,
		Body:  bytes.NewReader(body),
	}
	return c.exec(req)
}

// 分页查询，用于大数量查询，普通查询，默认最多返回10000条
func (c *ESClient) PageSearch(index []string, body []byte, fn func(p goo_utils.Params) error) error {
	var (
		scrollDuration = time.Minute
		size           = 500
		scrollId       string
	)
	defer func() {
		if scrollId == "" {
			return
		}
		res, clearErr := esapi.ClearScrollRequest{
			ScrollID: []string{scrollId},
		}.Do(context.Background(), c.cli)
		if res != nil && res.Body != nil {
			_ = res.Body.Close()
		}
		if clearErr != nil {
			c.log().Error(clearErr)
		}
	}()

	for n := 0; ; n++ {
		var (
			res *esapi.Response
			err error
		)

		if n == 0 {
			res, err = esapi.SearchRequest{
				Index:  index,
				Body:   bytes.NewReader(body),
				Scroll: scrollDuration,
				Size:   &size, // 每次获取的文档数量
			}.Do(context.Background(), c.Client())
		} else {
			res, err = esapi.ScrollRequest{
				ScrollID: scrollId,
				Scroll:   scrollDuration,
			}.Do(context.Background(), c.Client())
		}

		if err != nil {
			c.log().Error(err)
			return err
		}

		if res.IsError() {
			err = fmt.Errorf("error getting response: %s", res.String())
			_ = res.Body.Close()
			c.log().Error(err)
			return err
		}

		b, err := io.ReadAll(res.Body)
		_ = res.Body.Close()
		if err != nil {
			c.log().Error(err)
			return err
		}

		p, err := goo_utils.Byte(b).Params()
		if err != nil {
			c.log().Error(err)
			return err
		}

		if sid := p.Get("_scroll_id").String(); sid != "" {
			scrollId = sid
		}

		hits := p.Get("hits.hits").Array()
		if len(hits) == 0 {
			break
		}

		c.log().DebugF("第 %d 页，每页 %d 条", n+1, size)

		for _, hit := range hits {
			if err = fn(hit.Get("_source")); err != nil {
				c.log().Error(err)
				return err
			}
		}
	}

	return nil
}

// SearchAfter 用 search_after 深翻页，不受 scroll 超时影响。
// body 未指定 sort 时默认 [{"_id":"asc"}]；请勿与 from 同用（会忽略 from）。
func (c *ESClient) SearchAfter(index []string, body goo_utils.M, fn func(p goo_utils.Params) error) error {
	if body == nil {
		body = goo_utils.M{}
	}
	if _, ok := body["size"]; !ok {
		body["size"] = 500
	}
	if body["sort"] == nil {
		body["sort"] = []goo_utils.M{{"_id": "asc"}}
	}
	delete(body, "from")

	for n := 0; ; n++ {
		res, err := c.Search(index, body.Json())
		if err != nil {
			c.log().Error(err)
			return err
		}

		if res.IsError() {
			err = fmt.Errorf("error getting response: %s", res.String())
			_ = res.Body.Close()
			c.log().Error(err)
			return err
		}

		b, err := io.ReadAll(res.Body)
		_ = res.Body.Close()
		if err != nil {
			c.log().Error(err)
			return err
		}

		p, err := goo_utils.Byte(b).Params()
		if err != nil {
			c.log().Error(err)
			return err
		}

		hits := p.Get("hits.hits").Array()
		if len(hits) == 0 {
			return nil
		}

		c.log().DebugF("search_after 第 %d 页，%d 条", n+1, len(hits))

		for _, hit := range hits {
			if err = fn(hit.Get("_source")); err != nil {
				c.log().Error(err)
				return err
			}
		}

		sortVals := hits[len(hits)-1].Get("sort").ArrayData()
		if len(sortVals) == 0 {
			err = fmt.Errorf("search_after: missing hit.sort")
			c.log().Error(err)
			return err
		}
		body["search_after"] = sortVals
	}
}
