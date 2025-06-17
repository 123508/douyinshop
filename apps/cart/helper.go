package main

import (
	"github.com/123508/douyinshop/kitex_gen/cart"
	"github.com/123508/douyinshop/pkg/myredis"
	"github.com/123508/douyinshop/pkg/util"
	"github.com/redis/go-redis/v9"
	"strings"
)

const (
	serviceName = "cart"
)

var Rds = connectWithRedis()

func connectWithRedis() *redis.Client {
	rds, err := myredis.InitRedis()
	if err != nil {
		util.LogError("打开Redis连接失败", "connectWithRedis", "", err)
	}
	return rds
}

// 验证添加商品请求
func validateAddItemReq(req *cart.AddItemReq) error {
	if req == nil {
		return NilRequestError
	}
	if req.Item == nil {
		return NilItemError
	}
	if req.Item.Quantity <= 0 {
		return NegativeQuantityError
	}
	if req.UserId == 0 {
		return NilUserIdError
	}
	return nil
}

// IsDuplicateKeyError 检查是否为重复键错误
func IsDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "Duplicate entry")
}
