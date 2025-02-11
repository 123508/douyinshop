package main

import (
	userOrder "github.com/123508/douyinshop/kitex_gen/order/userOrder/orderuserservice"
	"log"
)

func main() {
	svr := userOrder.NewServer(new(OrderUserServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
