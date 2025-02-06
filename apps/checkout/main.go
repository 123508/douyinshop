package main

import (
	"log"

	checkout "github.com/123508/douyinshop/kitex_gen/checkout/checkoutservice"
)

func main() {
	svr := checkout.NewServer(new(CheckoutServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
