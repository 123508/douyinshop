package main

import (
	"fmt"
	payment "github.com/123508/douyinshop/kitex_gen/payment/paymentservice"
	"github.com/123508/douyinshop/pkg/config"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	"log"
	"net"
)

func main() {
	// etcd注册服务
	r, err := etcd.NewEtcdRegistryWithAuth(config.Conf.EtcdConfig.Endpoints, config.Conf.EtcdConfig.Username, config.Conf.EtcdConfig.Password)
	if err != nil {
		log.Fatal(err)
	}
	addr, _ := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", config.Conf.PaymentConfig.Host, config.Conf.PaymentConfig.Port))
	svr := payment.NewServer(
		new(PaymentServiceImpl),
		server.WithServiceAddr(addr),
		server.WithRegistry(r),
		server.WithServerBasicInfo(
			&rpcinfo.EndpointBasicInfo{
				ServiceName: config.Conf.PaymentConfig.ServiceName,
			},
		),
	)
	err = svr.Run()
	if err != nil {
		log.Println(err.Error())
	}
}
