package main

import (
	"github.com/123508/douyinshop/pkg/myredis"
	"gorm.io/gorm"
	"log"
	"net"
	"time"

	user "github.com/123508/douyinshop/kitex_gen/user/userservice"
	"github.com/123508/douyinshop/pkg/config"
	"github.com/123508/douyinshop/pkg/db"
	"github.com/123508/douyinshop/pkg/models"

	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"

	"fmt"
)

func main() {
	db, err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	rds, err := myredis.InitRedis()

	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.UserLogin{})
	db.AutoMigrate(&models.UserRole{})

	db.Model(&models.User{}).FirstOrCreate(&models.User{
		Model: gorm.Model{
			CreatedAt: time.Now(),
		},
		ID:     1,
		Name:   "管理员",
		Email:  "admin",
		Phone:  "admin",
		Gender: 1,
		Avatar: "",
		Status: 1,
	})

	db.Model(&models.UserLogin{}).FirstOrCreate(&models.UserLogin{
		Model: gorm.Model{
			CreatedAt: time.Now(),
		},
		ID:       1,
		UserId:   1,
		Password: "admin",
	})

	db.Model(&models.UserRole{}).FirstOrCreate(&models.UserRole{
		Model: gorm.Model{
			CreatedAt: time.Now(),
		},
		ID:     1,
		UserID: 1,
		RoleID: 1,
	})

	r, err := etcd.NewEtcdRegistryWithAuth(config.Conf.EtcdConfig.Endpoints, config.Conf.EtcdConfig.Username, config.Conf.EtcdConfig.Password)
	if err != nil {
		log.Fatal(err)
	}

	addr, _ := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", config.Conf.UserConfig.Host, config.Conf.UserConfig.Port))
	svr := user.NewServer(
		NewUserServiceImpl(db, rds),
		server.WithServiceAddr(addr),
		server.WithRegistry(r),
		server.WithServerBasicInfo(
			&rpcinfo.EndpointBasicInfo{
				ServiceName: config.Conf.UserConfig.ServiceName,
			},
		),
	)

	err = svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
