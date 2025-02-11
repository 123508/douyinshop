package main

import (
	address "github.com/123508/douyinshop/kitex_gen/address/addressservice"
	"log"
)

func main() {
	svr := address.NewServer(new(AddressServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
