package els

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/123508/douyinshop/pkg/config"
	"log"
	"strconv"

	"github.com/elastic/go-elasticsearch/v8"
)

var es *elasticsearch.Client

// 初始化ElasticSearch客户端
func init() {
	cfg := elasticsearch.Config{
		Addresses: config.Conf.ElasticSearch.Hosts,
		Username:  config.Conf.ElasticSearch.Username,
		Password:  config.Conf.ElasticSearch.Password,
	}
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the ElasticSearch client: %s", err)
	}

	es = client
}

// SearchProduct 搜索商品
// name: 商品名称
// 返回值: 商品id列表
func SearchProduct(name string, page int, size int) ([]uint32, error) {
	if es == nil {
		return nil, fmt.Errorf("ElasticSearch client is nil")
	}

	// 构建搜索请求
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"should": []map[string]interface{}{
					{
						"match": map[string]interface{}{
							"name": name,
						},
					},
					{
						"match": map[string]interface{}{
							"description": name,
						},
					},
					{
						"match": map[string]interface{}{
							"categories": name,
						},
					},
				},
			},
		},
		"from": (page - 1) * size,
		"size": size,
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, fmt.Errorf("error encoding query: %s", err)
	}
	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex("product"),
		es.Search.WithBody(&buf),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	)
	if err != nil {
		return nil, fmt.Errorf("error searching: %s", err)
	}
	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, fmt.Errorf("error parsing the response body: %s", err)
		}
		return nil, fmt.Errorf("error: %s", e)
	}

	// 获取res中所有搜索结果的id
	result := make([]uint32, 0)
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("error parsing the response body: %s", err)
	}
	hits := r["hits"].(map[string]interface{})["hits"].([]interface{})
	for _, hit := range hits {
		hit1 := hit.(map[string]interface{})
		id, _ := strconv.Atoi(hit1["_id"].(string))
		result = append(result, uint32(id))
	}
	return result, nil
}
