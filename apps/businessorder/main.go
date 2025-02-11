package main

import (
	businessOrder "github.com/123508/douyinshop/kitex_gen/order/businessOrder/orderbusinessservice"
	"log"
)

func main() {
	svr := businessOrder.NewServer(new(OrderBusinessServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
