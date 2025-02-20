package main

import (
	"fmt"
	"github.com/123508/douyinshop/pkg/config"
	"github.com/123508/douyinshop/pkg/db"
	"github.com/123508/douyinshop/pkg/models"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	"log"
	"net"

	checkout "github.com/123508/douyinshop/kitex_gen/checkout/checkoutservice"
)

func main() {
	database, err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	database.AutoMigrate(&models.AddressBook{})

	r, err := etcd.NewEtcdRegistryWithAuth(config.Conf.EtcdConfig.Endpoints, config.Conf.EtcdConfig.Username, config.Conf.EtcdConfig.Password)
	if err != nil {
		log.Fatal(err)
	}

	addr, _ := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", config.Conf.CheckoutConfig.Host, config.Conf.CheckoutConfig.Port))

	svr := checkout.NewServer(
		new(CheckoutServiceImpl),
		server.WithServiceAddr(addr),
		server.WithRegistry(r),
		server.WithServerBasicInfo(
			&rpcinfo.EndpointBasicInfo{
				ServiceName: config.Conf.CheckoutConfig.ServiceName,
			},
		),
	)

	err = svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
