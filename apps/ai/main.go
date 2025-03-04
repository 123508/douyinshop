package main

import (
    "fmt"
    "log"
    "net"

    ai "github.com/123508/douyinshop/kitex_gen/ai/aiservice"
    "github.com/123508/douyinshop/kitex_gen/order/userOrder/orderuserservice"
    "github.com/123508/douyinshop/pkg/config"
    "github.com/cloudwego/kitex/client"
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
        log.Fatalf("创建etcd注册器失败: %v", err)
    }

    // 创建订单服务客户端
    orderClient, err := orderuserservice.NewClient(
        config.Conf.OrderConfig.ServiceName,
        client.WithResolver(r),
    )
    if err != nil {
        log.Fatalf("创建订单服务客户端失败: %v", err)
    }

    // 创建服务地址
    addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", 
        config.Conf.AIConfig.Host, 
        config.Conf.AIConfig.Port))
    if err != nil {
        log.Fatalf("解析服务地址失败: %v", err)
    }

    // 创建服务实现实例
    impl, err := NewAiServiceImpl(orderClient)
    if err != nil {
        log.Fatalf("创建服务实现实例失败: %v", err)
    }
    defer impl.Close()

    // 创建服务实例
    svr := ai.NewServer(
        impl,
        server.WithServiceAddr(addr),
        server.WithRegistry(r),
        server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
            ServiceName: config.Conf.AIConfig.ServiceName,
        }),
    )

    // 启动服务
    if err := svr.Run(); err != nil {
        log.Fatalf("服务运行失败: %v", err)
    }
}