package main

import (
	"fmt"
	businessOrder "github.com/123508/douyinshop/kitex_gen/order/businessOrder/orderbusinessservice"
	"github.com/123508/douyinshop/pkg/config"
	"github.com/123508/douyinshop/pkg/db"
	"github.com/123508/douyinshop/pkg/models"
	"github.com/123508/douyinshop/pkg/util"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	"github.com/streadway/amqp"
	"log"
	"net"
	"sync"
)

type RabbitmqHandler struct {
}

func (t RabbitmqHandler) HandlerMessage(msg amqp.Delivery) {
	fmt.Println(msg.Body)
}

func main() {

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		util.ReceiveMessages("order.direct", util.Direct, "order.queue", "message", &RabbitmqHandler{})

	}()

	go func() {
		defer wg.Done()
		// 初始化数据库
		db, err := db.InitDB()
		if err != nil {
			log.Fatal(err)
		}
		db.AutoMigrate(&models.Order{})
		db.AutoMigrate(&models.OrderDetail{})
		db.AutoMigrate(&models.OrderStatusLog{})
		db.AutoMigrate(&models.Product{})

		// 设置 Etcd 注册
		r, err := etcd.NewEtcdRegistryWithAuth(config.Conf.EtcdConfig.Endpoints, config.Conf.EtcdConfig.Username, config.Conf.EtcdConfig.Password)
		if err != nil {
			log.Fatal(err)
		}

		// 设置服务地址
		addr, _ := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", config.Conf.BusinessOrderConfig.Host, config.Conf.BusinessOrderConfig.Port))

		// 创建 server
		svr := businessOrder.NewServer(
			new(OrderBusinessServiceImpl),
			server.WithServiceAddr(addr),
			server.WithRegistry(r),
			server.WithServerBasicInfo(
				&rpcinfo.EndpointBasicInfo{
					ServiceName: config.Conf.BusinessOrderConfig.ServiceName,
				},
			),
		)

		// 启动 server
		err = svr.Run()
		if err != nil {
			log.Println(err.Error())
		}
	}()

	wg.Wait()

}
