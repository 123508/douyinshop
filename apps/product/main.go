package main

import (
	"fmt"
	product "github.com/123508/douyinshop/kitex_gen/product/productcatalogservice"
	"github.com/123508/douyinshop/pkg/config"
	"github.com/123508/douyinshop/pkg/db"
	"github.com/123508/douyinshop/pkg/models"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	"gorm.io/gorm"
	"log"
	"net"
)

var database *gorm.DB

func main() {
	databaseConnection, err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	database = databaseConnection

	err = database.AutoMigrate(&models.Product{})
	if err != nil {
		log.Fatal(err)
	}
	err = database.AutoMigrate(&models.Category{})
	if err != nil {
		log.Fatal(err)
	}

	r, err := etcd.NewEtcdRegistryWithAuth(config.Conf.EtcdConfig.Endpoints, config.Conf.EtcdConfig.Username, config.Conf.EtcdConfig.Password)
	if err != nil {
		log.Fatal(err)
	}

	addr, _ := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", config.Conf.ProductConfig.Host, config.Conf.ProductConfig.Port))
	svr := product.NewServer(
		new(ProductCatalogServiceImpl),
		server.WithServiceAddr(addr),
		server.WithRegistry(r),
		server.WithServerBasicInfo(
			&rpcinfo.EndpointBasicInfo{
				ServiceName: config.Conf.ProductConfig.ServiceName,
			},
		),
	)

	err = svr.Run()
	if err != nil {
		log.Println(err.Error())
	}
}
