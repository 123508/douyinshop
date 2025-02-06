package main

import (
	"log"
	"net"

	user "github.com/123508/douyinshop/kitex_gen/user/userservice"
	"github.com/123508/douyinshop/pkg/config"
	"github.com/123508/douyinshop/pkg/db"
	"github.com/123508/douyinshop/pkg/models"

	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"

	"fmt"
)

func main() {
	db, err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&models.User{})

	r, err := etcd.NewEtcdRegistryWithAuth(config.Conf.EtcdConfig.Endpoints, config.Conf.EtcdConfig.Username, config.Conf.EtcdConfig.Password)
	if err != nil {
		log.Fatal(err)
	}

	addr, _ := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", config.Conf.UserConfig.Host, config.Conf.UserConfig.Port))
	svr := user.NewServer(
		new(UserServiceImpl),
		server.WithServiceAddr(addr),
		server.WithRegistry(r),
		server.WithServerBasicInfo(
			&rpcinfo.EndpointBasicInfo{
				ServiceName: config.Conf.UserConfig.ServiceName,
			},
		),
	)

	err = svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
