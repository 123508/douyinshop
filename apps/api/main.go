package main

import (
	"log"

	"github.com/123508/douyinshop/apps/api/handlers/auth"
	"github.com/123508/douyinshop/apps/api/handlers/user"
	"github.com/123508/douyinshop/pkg/config"

	"fmt"

	"github.com/cloudwego/hertz/pkg/app/server"
)

func main() {
	hertzAddr := fmt.Sprintf("%s:%d", config.Conf.HertzConfig.Host, config.Conf.HertzConfig.Port)
	hz := server.New(server.WithHostPorts(hertzAddr))

	userGrop := hz.Group("/user")
	userGrop.GET("/register", user.Register)
	userGrop.GET("/login", user.Login)

	authGroup := hz.Group("/auth")
	authGroup.GET("/deliver_token", auth.DeliverToken)
	authGroup.GET("/verify_token", auth.VerifyToken)

	if err := hz.Run(); err != nil {
		log.Fatal(err)
	}
}
