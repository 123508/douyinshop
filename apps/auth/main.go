package main

import (
	"context"
	"log"
	"net"
	"time"

	auth "github.com/123508/douyinshop/kitex_gen/auth/authservice"
	"github.com/123508/douyinshop/pkg/config"
	"github.com/123508/douyinshop/pkg/redis"

	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"

	"fmt"
)

func main() {
	rdb, err := redis.InitRedis()
	if err != nil {
		log.Fatal(err)
	}
	rdb.Set(context.Background(), "key", "value", 60*time.Second)

	r, err := etcd.NewEtcdRegistryWithAuth(config.Conf.EtcdConfig.Endpoints, config.Conf.EtcdConfig.Username, config.Conf.EtcdConfig.Password)
	if err != nil {
		log.Fatal(err)
	}

	addr, _ := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", config.Conf.AuthConfig.Host, config.Conf.AuthConfig.Port))
	svr := auth.NewServer(
		new(AuthServiceImpl),
		server.WithServiceAddr(addr),
		server.WithRegistry(r),
		server.WithServerBasicInfo(
			&rpcinfo.EndpointBasicInfo{
				ServiceName: config.Conf.AuthConfig.ServiceName,
			},
		),
	)

	err = svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
