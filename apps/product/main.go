package main

import (
	"log"

	product "github.com/123508/douyinshop/kitex_gen/product/productcatalogservice"
)

func main() {
	svr := product.NewServer(new(ProductCatalogServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
