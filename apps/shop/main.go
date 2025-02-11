package main

import (
	shop "github.com/123508/douyinshop/kitex_gen/shop/shopservice"
	"log"
)

func main() {
	svr := shop.NewServer(new(ShopServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
