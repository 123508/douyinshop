# 抖音电商

### 项目运行

#### 运行所需环境

- Go 1.23.1
- MySQL 8.0
- Redis 6.0
- etcd 3.5

可通过docker-compose启动所需环境

```shell
docker-compose up -d
```

#### 配置文件

> 配置文件位于`config`目录下，根据实际情况修改配置文件

- mysql: mysql相关配置
- redis: redis相关配置
- etcd: etcd相关配置
- elasticsearch: elasticsearch相关配置

#### 服务启动

api网关: api

```shell
cd apps/api
go run .
```

认证中心: auth

```shell
cd apps/auth
go run .
```

购物车服务: cart

```shell
cd apps/cart
go run .
```

结算服务: checkout

```shell
cd apps/checkout
go run .
```

订单服务: order

```shell
cd apps/order
go run .
```

支付服务: payment

```shell
cd apps/payment
go run .
```

商品服务: product

```shell
cd apps/product
go run .
```

用户服务: user

```shell
cd apps/user
go run .
```

地址服务: address

```shell
cd apps/address
go run .
```

商家订单服务: businessorder

```shell
cd apps/businessorder
go run .
```

店铺服务: shop

```shell
cd apps/shop
go run .
```

AI服务: ai

```shell
cd apps/ai
go run .
```

### 项目结构

```tree
.
├── apps                          // 模块逻辑目录
│   ├── api                       // api网关
│   ├── auth                      // 认证中心
│   ├── cart                      // 购物车服务
│   ├── checkout                  // 结算服务
│   ├── order                     // 订单服务
│   ├── payment                   // 支付服务
│   ├── product                   // 商品服务
│   ├── address                   // 地址服务
│   ├── businessorder             // 商家订单服务
│   ├── shop                      // 店铺服务
│   ├── ai                        // AI服务
│   └── user                      // 用户服务
├── config                        // 配置文件目录
├── deploy                        // 项目所需的环境部署
├── docker-compose.yml            // docker-compose文件
├── kitex_gen                     // Kitex生成的代码
├── kitex_gen.sh                  // kitex生成脚手架代码脚本
├── pkg                           // 项目所依赖的一些公共包
│   ├── ai                        // AI相关包
│   ├── config                    // 读取配置文件
│   ├── db                        // 数据库初始化
│   ├── models                    // 基本模型定义
│   ├── els                       // ElasticSearch初始化
│   ├── errors                    // 错误定义
│   ├── utils                     // 工具包
│   └── redis                     // redis初始化
├── proto                         // Protobuf
└── README.md
```