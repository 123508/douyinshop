package main

import (
	"fmt"
	"github.com/123508/douyinshop/pkg/config"
	"log"

	product "github.com/123508/douyinshop/kitex_gen/product/productcatalogservice"
	"github.com/elastic/go-elasticsearch/v8"
)

func main() {
	cfg := elasticsearch.Config{
		Addresses: config.Conf.ElasticSearch.Hosts,
		Username:  config.Conf.ElasticSearch.Username,
		Password:  config.Conf.ElasticSearch.Password,
	}
	els, err := elasticsearch.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	// 测试代码
	x, err := els.Get("product", "1")
	if err != nil {
		panic(err)
	}
	fmt.Println(x)

	svr := product.NewServer(new(ProductCatalogServiceImpl))

	err = svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
