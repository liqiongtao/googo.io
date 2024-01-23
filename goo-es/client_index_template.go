package goo_es

import (
	"bytes"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

// 删除索引模版
func (c *client) IndexTemplateDel(name string) (*esapi.Response, error) {
	req := esapi.IndicesDeleteTemplateRequest{
		Name: name,
	}
	return c.exec(req)
}

// 创建索引模版
//
// - lifecycle.name: 生命周期策略
// - number_of_shards: 索引分片数量
// - number_of_replicas: 索引副本数量
// - - nested: 嵌套对象 {list: [{id:1, name:''}]}
//
// {
//    "index_patterns": ["abc-*"],
//    "template": {
//        "settings": {
//            "number_of_shards": 3,
//            "number_of_replicas": 0,
//            "index.lifecycle.name": "90-days-default",
//            "index.lifecycle.rollover_alias": "abc"
//        },
//        "mappings": {
//            "properties": {
//                "id": {"type": "integer"},
//                "datetime": {"type": "date", "format": "yyyy-MM-dd HH:mm:ss"},
//                "name": {"type": "text"},
//                "list": {
//                    "type": "nested",
//                    "properties": {
//                        "id": {"type": "integer"},
//                        "name": {"type": "text"},
//                    }
//                }
//            }
//        }
//    }
// }
func (c *client) IndexTemplatePut(name string, b []byte) (*esapi.Response, error) {
	req := esapi.IndicesPutIndexTemplateRequest{
		Name: name,
		Body: bytes.NewReader(b),
	}
	return c.exec(req)
}
