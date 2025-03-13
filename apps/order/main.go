package main

import (
	"fmt"
	order "github.com/123508/douyinshop/kitex_gen/order/userOrder/orderuserservice"
	"github.com/123508/douyinshop/pkg/config"
	"github.com/123508/douyinshop/pkg/db"
	"github.com/123508/douyinshop/pkg/models"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
	"log"
	"net"
	"sync"
	"time"
)

func main() {

	//无缓冲通道阻塞队列
	DBChannel := make(chan *gorm.DB)

	stopChannel := make(chan bool, 1)

	exitChannel := make(chan bool, 1)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		// 创建一个新的cron实例

		select {

		case <-stopChannel:
			log.Println("任务调度器退出")
			return

		default:

			DB := <-DBChannel

			//检查数据库是否初始化成功
			if DB.Error != nil {
				log.Println("数据库初始化出错")
				stopChannel <- true
				return
			}

			log.Println("任务调度器初始化完毕")

			c := cron.New(cron.WithSeconds(), // 使用WithSeconds()来支持秒级别的调度
				cron.WithLocation(time.Local), // 使用本地时区
			)

			defer c.Stop()

			// 添加一个每分钟执行一次的任务
			_, err := c.AddFunc("0 */1 * * * ?", func() {
				log.Println("当前时间:", time.Now())
				var logs []models.OrderStatusLog

				if err := DB.Model(&models.OrderStatusLog{}).Where("status = 0 and end_time = null").Find(&logs).Error; err != nil {
					log.Println(err)
					return
				}

				err := DB.Transaction(func(tx *gorm.DB) error {

					for _, orderly := range logs {

						current := time.Now()
						//超时订单
						if current.After(orderly.StartTime.Add(15 * time.Minute)) {
							//先修改当前时间
							if err := DB.Where("id = ?", orderly.ID).Update("end_time", &current).Error; err != nil {
								return err
							}
							//准备插入新节点
							newLog := models.OrderStatusLog{
								OrderId:     orderly.OrderId,
								Status:      6,
								StartTime:   &current,
								EndTime:     nil,
								Description: "订单超时,已取消",
							}
							//更新新节点为已取消状态
							if err := DB.Create(&newLog).Error; err != nil {
								return err
							}
							//向order中写入状态
							if err := DB.Model(&models.Order{}).Where("id = ?", orderly.OrderId).Update("status", 6).Error; err != nil {
								return err
							}
						}
					}

					return nil
				})

				if err != nil {
					log.Println(err)
					return
				}
			})

			if err != nil {
				log.Println(err)
				stopChannel <- true
				return
			}

			log.Println("开始执行任务")

			// 启动cron调度器
			c.Start()

			//阻塞监听(必须阻塞,否则任务调度器会自动退出)
			for {
				if <-stopChannel {
					stopChannel <- true
					break
				}
			}
		}
	}()

	go func() {
		defer wg.Done()

		select {
		case <-exitChannel:
			log.Println("主协程退出")
			return
		default:

			// 初始化数据库
			db, err := db.InitDB()
			if err != nil {
				log.Println(err)
				stopChannel <- true
				exitChannel <- true
				return
			}
			//将初始化数据传递给另一个协程
			DBChannel <- db

			db.AutoMigrate(&models.Order{})
			db.AutoMigrate(&models.OrderDetail{})
			db.AutoMigrate(&models.OrderStatusLog{})

			// 设置 Etcd 注册
			r, err := etcd.NewEtcdRegistryWithAuth(config.Conf.EtcdConfig.Endpoints, config.Conf.EtcdConfig.Username, config.Conf.EtcdConfig.Password)
			if err != nil {
				log.Println(err)
				stopChannel <- true
				exitChannel <- true
				return
			}

			// 设置服务地址
			addr, _ := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", config.Conf.OrderConfig.Host, config.Conf.OrderConfig.Port))

			// 创建 server
			svr := order.NewServer(
				new(OrderUserServiceImpl), // 创建订单服务实例
				server.WithServiceAddr(addr),
				server.WithRegistry(r),
				server.WithServerBasicInfo(
					&rpcinfo.EndpointBasicInfo{
						ServiceName: config.Conf.OrderConfig.ServiceName, // 设置服务名
					},
				),
			)

			// 启动 server
			err = svr.Run()
			if err != nil {
				log.Println(err.Error())
				stopChannel <- true
				exitChannel <- true
				return
			}
		}
	}()

	wg.Wait()
}
