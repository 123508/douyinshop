package main

import (
	ai "github.com/123508/douyinshop/kitex_gen/ai/aiservice"
	"log"
)

func main() {
	svr := ai.NewServer(new(AiServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
