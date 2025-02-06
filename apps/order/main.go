package main

import (
	"log"

	order "github.com/123508/douyinshop/kitex_gen/order/orderservice"
)

func main() {
	svr := order.NewServer(new(OrderServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
