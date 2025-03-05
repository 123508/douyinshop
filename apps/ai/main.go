package main

import (
    "fmt"
    "log"
    "net"

    ai "github.com/123508/douyinshop/kitex_gen/ai/aiservice"
    "github.com/123508/douyinshop/kitex_gen/order/userOrder/orderuserservice"
    "github.com/123508/douyinshop/pkg/config"
    "github.com/cloudwego/kitex/client"
    "github.com/cloudwego/kitex/pkg/discovery"
    "github.com/cloudwego/kitex/pkg/rpcinfo"
    "github.com/cloudwego/kitex/server"
    etcd "github.com/kitex-contrib/registry-etcd"
)

func main() {
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
    orderClient, err := orderuserservice.NewClient(
        config.Conf.OrderConfig.ServiceName,
        client.WithResolver(r.(discovery.Resolver)),
    )
    if err != nil {
        log.Fatal(err)
    }

    // 创建服务地址
    addr, _ := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", 
        config.Conf.AIConfig.Host, 
        config.Conf.AIConfig.Port))

    // 创建服务实例
    svr := ai.NewServer(
        NewAiServiceImpl(orderClient),
        server.WithServiceAddr(addr),
        server.WithRegistry(r),
        server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
            ServiceName: config.Conf.AIConfig.ServiceName,
        }),
    )

    // 启动服务
    err = svr.Run()
    if err != nil {
        log.Println(err.Error())
    }
}