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
		ctx            = context.Background()
		scrollDuration = 2 * time.Second
		size           = 500
		scrollId       string
	)

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
			}.Do(ctx, c.Client())
		} else {
			res, err = esapi.ScrollRequest{
				ScrollID: scrollId,
				Scroll:   scrollDuration,
			}.Do(ctx, c.Client())
		}

		if err != nil {
			c.log().Error(err)
			return err
		}

		if res.IsError() {
			c.log().Error(fmt.Errorf("error getting initial response: %s", res.String()))
			return fmt.Errorf("error getting initial response: %s", res.String())
		}

		// 获取数据
		b, err := io.ReadAll(res.Body)
		if err != nil {
			c.log().Error(err)
			return err
		}

		// 关闭数据流
		res.Body.Close()

		// 转换
		p, err := goo_utils.Byte(b).Params()
		if err != nil {
			c.log().Error(err)
			return err
		}

		if len(p.Get("hits.hits").Array()) == 0 {
			break
		}

		// 滚动ID
		scrollId = p.Get("_scroll_id").String()

		// 打印翻页
		c.log().DebugF("第 %d 页，每页 %d 条", n+1, size)

		// 处理文档
		for _, hit := range p.Get("hits.hits").Array() {
			if err = fn(hit.Get("_source")); err != nil {
				c.log().Error(err)
				return nil
			}
		}
	}

	// 清理Scroll上下文
	_, err := esapi.ClearScrollRequest{
		ScrollID: []string{scrollId},
	}.Do(context.Background(), c.cli)
	if err != nil {
		c.log().Error(err)
	}

	return nil
}
