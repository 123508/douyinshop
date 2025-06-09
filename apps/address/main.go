package main

import (
	"fmt"
	address "github.com/123508/douyinshop/kitex_gen/address/addressservice"
	"github.com/123508/douyinshop/pkg/config"
	"github.com/123508/douyinshop/pkg/db"
	"github.com/123508/douyinshop/pkg/models"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	log "github.com/sirupsen/logrus"
	"net"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		FieldMap: log.FieldMap{
			log.FieldKeyTime:  "时间",
			log.FieldKeyLevel: "日志类型",
			log.FieldKeyMsg:   "日志内容",
		},
	})
}

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

	addr, _ := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", config.Conf.AddressConfig.Host, config.Conf.AddressConfig.Port))

	svr := address.NewServer(
		new(AddressServiceImpl),
		server.WithServiceAddr(addr),
		server.WithRegistry(r),
		server.WithServerBasicInfo(
			&rpcinfo.EndpointBasicInfo{
				ServiceName: config.Conf.AddressConfig.ServiceName,
			},
		),
	)

	err = svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
