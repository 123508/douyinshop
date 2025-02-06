package main

import (
	"log"

	cart "github.com/123508/douyinshop/kitex_gen/cart/cartservice"
)

func main() {
	svr := cart.NewServer(new(CartServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
