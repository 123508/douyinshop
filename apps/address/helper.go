package main

import (
	"context"
	"github.com/123508/douyinshop/kitex_gen/address"
	"github.com/123508/douyinshop/pkg/db"
	"github.com/123508/douyinshop/pkg/models"
	"github.com/123508/douyinshop/pkg/myredis"
	"github.com/123508/douyinshop/pkg/util"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const (
	serviceName = "address"
	info        = "info"
)

var DB = connectWithMySQL()

func connectWithMySQL() *gorm.DB {
	DB, err := db.InitDB()
	if err != nil {
		util.LogError("打开MySQL连接失败", "connectWithMySQL", "", err)
	}
	return DB
}

var Rds = connectWithRedis()

func connectWithRedis() *redis.Client {
	rds, err := myredis.InitRedis()
	if err != nil {
		util.LogError("打开Redis连接失败", "connectWithRedis", "", err)
	}
	return rds
}

// 两个地址转换函数,注意地址类型有Address,AddressItem,AddressBook
func tranAddressToAddressBook(origin *address.Address) *models.AddressBook {
	addr := &models.AddressBook{}
	addr.StressAddress = origin.StreetAddress
	addr.Phone = origin.Phone
	addr.Gender = origin.Gender
	addr.Consignee = origin.Consignee
	addr.State = origin.State
	addr.City = origin.City
	addr.Country = origin.Country
	addr.Label = origin.Label
	addr.ZipCode = origin.ZipCode //共计九个字段
	return addr
}

func tranAddressBookToAddress(origin *models.AddressBook) *address.Address {
	addr := &address.Address{}
	addr.StreetAddress = origin.StressAddress
	addr.Phone = origin.Phone
	addr.Gender = origin.Gender
	addr.Consignee = origin.Consignee
	addr.State = origin.State
	addr.City = origin.City
	addr.Country = origin.Country
	addr.Label = origin.Label
	addr.ZipCode = origin.ZipCode
	addr.IsDefault = origin.IsDefault //共计十个字段
	return addr
}

func DelUserOrderCache(ctx context.Context, rds *redis.Client, userId interface{}) {
	// 构建 pattern
	pattern := util.TakeKey(serviceName, userId, "*")
	var cursor uint64 = 0
	var batchSize int64 = 100 // 每次扫描的数量，可根据实际情况调整

	for {
		// 使用 SCAN 命令遍历所有匹配的 key
		keys, nextCursor, err := rds.Scan(ctx, cursor, pattern, batchSize).Result()
		if err != nil {
			util.LogError("redis扫描错误", "DelUserOrderCache", "", err)
		}
		if len(keys) > 0 {
			// pipeline 批量删除，提升性能
			pipe := rds.Pipeline()
			for _, key := range keys {
				pipe.Del(ctx, key)
			}
			if _, err := pipe.Exec(ctx); err != nil {
				util.LogError("redis pipeline删除错误", "DelUserOrderCache", "", err)
			}
		}
		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}
}
