package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/123508/douyinshop/pkg/db"
	"github.com/123508/douyinshop/pkg/models"
	"github.com/123508/douyinshop/pkg/myredis"
	"github.com/123508/douyinshop/pkg/util"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

const (
	goods       = "goods"
	serviceName = "businessOrder"
)

//订单列表的缓存方案 : id列表+详情分级缓存(后续应用到order/cart和product处)

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

func cleanCache(ctx context.Context, orderId uint64) {
	//清除缓存
	util.CleanCache(Rds, ctx, util.TakeKey(goods, "order", orderId))
	util.CleanCache(Rds, ctx, util.TakeKey(goods, "orderDetail", orderId))
	util.CleanCache(Rds, ctx, util.TakeKey(goods, "orderLog", orderId))
}

// GetOrderInfo 查询订单信息
func GetOrderInfo(ctx context.Context, orderId uint64) (*models.Order, error) {
	// 查询订单信息（Order）基本信息
	var order models.Order

	// 如果没有找到对应订单，返回错误信息
	if err := DB.Where("id = ?", orderId).First(&order).Error; err != nil {
		log.Println(err)
		return nil, SearchOrderError
	}

	return &order, nil
}

func GetOrderInfoWithCache(ctx context.Context, orderId uint64) (*models.Order, error) {
	// 查询订单信息（Order）基本信息
	var order models.Order

	orderKey := util.TakeKey(goods, "order", orderId)

	//查询单条order缓存
	jsonData, err := Rds.Get(ctx, orderKey).Result()

	if errors.Is(err, redis.Nil) {
		log.Println("order缓存未命中")
	} else if err != nil {
		util.LogError("查询order缓存失败", "GetOrderInfoWithCache", "查询order缓存", err)
	}

	//查询缓存失败,直接去查数据库
	if err != nil || json.Unmarshal([]byte(jsonData), &order) != nil {
		// 如果没有找到对应订单，返回错误信息
		if err := DB.Where("id = ?", orderId).First(&order).Error; err != nil {
			log.Println(err)
			return nil, SearchOrderError
		}

		jsonData, err := json.Marshal(&order)

		if err != nil {
			util.LogError("序列化订单数据错误", "GetOrderInfoWithCache", "", err)
		} else {
			//将数据加入缓存
			if setErr := Rds.Set(ctx, orderKey, jsonData, time.Duration(rand.Intn(3)+3)*time.Minute).Err(); setErr != nil {
				util.LogError("缓存数据失败", "GetOrderInfoWithCache", "将order加入缓存", setErr)
			}
		}
	} else {
		log.WithFields(log.Fields{
			"方法名": "GetOrderInfoWithCache",
		}).Info("查询order缓存成功")
	}

	return &order, nil
}

// GetOrderDetailsWithCache 查询指定订单的详细信息
func GetOrderDetailsWithCache(ctx context.Context, orderId uint64) ([]models.OrderDetail, error) {
	// 查询订单详情（List）
	orderDetails := make([]models.OrderDetail, 0)

	detailKey := util.TakeKey(goods, "orderDetail", orderId)

	//查询orderDetails缓存
	result, err := Rds.Get(ctx, detailKey).Result()

	//查询orderDetail缓存失败
	if err != nil || json.Unmarshal([]byte(result), &orderDetails) != nil {

		// 如果查询订单详情失败，返回错误
		if dbErr := DB.Where("order_id = ?", orderId).Find(&orderDetails).Error; dbErr != nil {
			util.LogError("查询订单详情失败", "GetOrderDetailsWithCache", "查询订单详情", dbErr)
			return nil, SearchOrderDetailsError
		}

		jsonData, err := json.Marshal(&orderDetails)

		if err != nil {
			util.LogError("序列化数据失败", "GetOrderDetailsWithCache", "序列化orderDetail", err)
		} else {
			if setErr := Rds.Set(ctx, detailKey, jsonData, time.Duration(rand.Intn(3)+3)*time.Minute).Err(); setErr != nil {
				util.LogError("缓存数据失败", "GetOrderDetailsWithCache", "将orderDetail加入缓存", setErr)
			}
		}

	} else {
		log.WithFields(log.Fields{
			"方法名": "GetOrderDetailsWithCache",
		}).Info("查询orderDetail缓存成功")
	}

	return orderDetails, nil
}

// GetOrderLogsWithCache 查询指定订单的日志信息
func GetOrderLogsWithCache(ctx context.Context, orderId uint64) ([]models.OrderStatusLog, error) {
	//查询订单日志详情
	orderLogs := make([]models.OrderStatusLog, 0)

	logKey := util.TakeKey(goods, "orderLog", orderId)

	res, err := Rds.Get(ctx, logKey).Result()

	//查询缓存失败
	if err != nil || json.Unmarshal([]byte(res), &orderLogs) != nil {
		//查询失败报错
		if err = DB.Where("order_id = ?", orderId).Find(&orderLogs).Error; err != nil {
			log.Println(err)
			return nil, SearchOrderLogsError
		}

		jsonData, err := json.Marshal(&orderLogs)

		if err != nil {
			util.LogError("序列化数据失败", "GetOrderLogsWithCache", "序列化orderLog", err)
		} else {
			if setErr := Rds.Set(ctx, logKey, jsonData, time.Duration(rand.Intn(3)+3)*time.Minute).Err(); setErr != nil {
				util.LogError("缓存数据失败", "GetOrderLogsWithCache", "将orderLog加入缓存", setErr)
			}

		}
	} else {
		log.WithFields(log.Fields{
			"方法名": "GetOrderLogsWithCache",
		}).Info("查询orderLog缓存成功")
	}

	return orderLogs, nil
}
