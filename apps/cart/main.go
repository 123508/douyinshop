package main

import (
    "fmt"
    "log"
    "net"

    "github.com/123508/douyinshop/kitex_gen/cart/cartservice"
    "github.com/123508/douyinshop/pkg/config"
    "github.com/123508/douyinshop/pkg/db"
    "github.com/123508/douyinshop/pkg/models"
    "github.com/cloudwego/kitex/pkg/rpcinfo"
    "github.com/cloudwego/kitex/server"
    etcd "github.com/kitex-contrib/registry-etcd"
)

func main() {
    // 初始化数据库
    database, err := db.InitDB()
    if err != nil {
        log.Fatal("Failed to initialize database:", err)
    }

    // 获取数据库实例，用于后续关闭
    sqlDB, err := database.DB()
    if err != nil {
        log.Fatal("Failed to get database instance:", err)
    }
    // 程序结束时关闭数据库连接
    defer sqlDB.Close()

    // 自动迁移 Cart 模型
    if err := database.AutoMigrate(&models.Cart{}); err != nil {
        log.Fatal("Failed to migrate cart model:", err)
    }

    // 创建ETCD注册中心
    r, err := etcd.NewEtcdRegistryWithAuth(
        config.Conf.EtcdConfig.Endpoints,
        config.Conf.EtcdConfig.Username,
        config.Conf.EtcdConfig.Password,
    )
    if err != nil {
        log.Fatal("Failed to connect to etcd:", err)
    }

    // 解析服务地址
    serviceAddr := fmt.Sprintf("%s:%d", config.Conf.CartConfig.Host, config.Conf.CartConfig.Port)
    addr, err := net.ResolveTCPAddr("tcp", serviceAddr)
    if err != nil {
        log.Fatal("Failed to resolve address:", err)
    }

    // 创建Kitex服务实例
    svr := cartservice.NewServer(
        &CartServiceImpl{db: database},
        server.WithServiceAddr(addr),
        server.WithRegistry(r),
        server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
            ServiceName: config.Conf.CartConfig.ServiceName,
        }),
    )

    // 启动服务
    log.Printf("Cart service starting on %s", serviceAddr)
    if err := svr.Run(); err != nil {
        log.Fatal("Service shutdown with error:", err) // 改用 Fatal，因为这是致命错误
    }
}