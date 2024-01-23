package goo_es

import (
	"bytes"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

// 添加、更新文档
func (c *ESClient) Index(index, docId string, b []byte) (*esapi.Response, error) {
	req := esapi.IndexRequest{
		Index:      index,
		DocumentID: docId,
		Body:       bytes.NewReader(b),
		Refresh:    "true",
	}
	return c.exec(req)
}

// 删除文档
func (c *ESClient) Delete(index, docId string) (*esapi.Response, error) {
	req := esapi.DeleteRequest{
		Index:      index,
		DocumentID: docId,
		Refresh:    "true",
	}
	return c.exec(req)
}

// 删除文档
func (c *ESClient) DeleteByQuery(index []string, b []byte) (*esapi.Response, error) {
	req := esapi.DeleteByQueryRequest{
		Index: index,
		Body:  bytes.NewReader(b),
	}
	*req.Refresh = true
	return c.exec(req)
}
