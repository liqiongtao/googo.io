package goo_es

import (
	"bytes"
	"github.com/elastic/go-elasticsearch/v7/esapi"
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
// {
//   "from" : 0,
//   "size" : 10,
//   "query": {
//     "bool": {
//       "must": [
//         {"match": { "user_id": 2716 }},
//         {"range": { "datetime": { "gte": "2024-01-14 00:00:00", "lte": "2024-01-14 23:59:59" }}},
//         {"range": { "size": { "gte": 60, "lte": 3600 }}},
//         {
//           "nested": {
//             "path": "list",
//             "query": {
//               "match_phrase": { "list.title": "青岛" }
//             },
//             "inner_hits": {
//               "size": 1
//             }
//           }
//         }
//       ]
//     }
//   },
//   "highlight": {
//     "pre_tags": "<em>",
//     "post_tags": "</em>",
//     "fields": {
//       "list.title": {}
//     }
//   },
//   "_source": {
//     "excludes" : ["list"]
//   }
// }
func (c *client) Search(index []string, b []byte) (*esapi.Response, error) {
	req := esapi.SearchRequest{
		Index: index,
		Body:  bytes.NewReader(b),
	}
	return c.exec(req)
}
