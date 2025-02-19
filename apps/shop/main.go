package main

import (
	"fmt"
	"gorm.io/gorm"
	"log"
	"net"

	"github.com/123508/douyinshop/kitex_gen/shop/shopservice"
	"github.com/123508/douyinshop/pkg/config"
	"github.com/123508/douyinshop/pkg/db"
	"github.com/123508/douyinshop/pkg/models"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
)

type ShopServiceImpl struct {
	db *gorm.DB
}

func main() {
	// 初始化数据库
	database, err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	database.AutoMigrate(&models.Shop{}, &models.Product{})

	// 创建ETCD注册中心
	r, err := etcd.NewEtcdRegistryWithAuth(
		config.Conf.EtcdConfig.Endpoints,
		config.Conf.EtcdConfig.Username,
		config.Conf.EtcdConfig.Password,
	)
	if err != nil {
		log.Fatal(err)
	}

	// 解析服务地址
	addr, _ := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d",
		config.Conf.ShopConfig.Host,
		config.Conf.ShopConfig.Port,
	))

	// 创建Kitex服务实例
	svr := shopservice.NewServer(
		&ShopServiceImpl{db: database},
		server.WithServiceAddr(addr),
		server.WithRegistry(r),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: config.Conf.ShopConfig.ServiceName,
		}),
	)

	// 启动服务
	if err := svr.Run(); err != nil {
		log.Println("Service shutdown with error:", err)
	}
}
