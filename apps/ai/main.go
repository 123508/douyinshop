package main

import (
    "fmt"
    "log"
    "net"

    ai "github.com/123508/douyinshop/kitex_gen/ai/aiservice"
    "github.com/123508/douyinshop/kitex_gen/order/orderservice"
    "github.com/123508/douyinshop/pkg/config"
    "github.com/123508/douyinshop/pkg/db"
    "github.com/cloudwego/kitex/client"
    "github.com/cloudwego/kitex/pkg/rpcinfo"
    "github.com/cloudwego/kitex/server"
    etcd "github.com/kitex-contrib/registry-etcd"
)

func main() {
    // 初始化数据库连接
    database, err := db.InitDB()
    if err != nil {
        log.Fatal(err)
    }

    // 创建 etcd 注册器
    r, err := etcd.NewEtcdRegistryWithAuth(
        config.Conf.EtcdConfig.Endpoints,
        config.Conf.EtcdConfig.Username,
        config.Conf.EtcdConfig.Password,
    )
    if err != nil {
        log.Fatal(err)
    }

    // 创建订单服务客户端
    orderClient, err := orderservice.NewClient(
        config.Conf.OrderConfig.ServiceName,
        client.WithResolver(r),
    )
    if err != nil {
        log.Fatal(err)
    }

    // 创建服务地址
    addr, _ := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", 
        config.Conf.AiConfig.Host, 
        config.Conf.AiConfig.Port))

    // 创建服务实例
    svr := ai.NewServer(
        NewAiServiceImpl(orderClient),
        server.WithServiceAddr(addr),
        server.WithRegistry(r),
        server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
            ServiceName: config.Conf.AiConfig.ServiceName,
        }),
    )

    // 启动服务
    err = svr.Run()
    if err != nil {
        log.Println(err.Error())
    }
}