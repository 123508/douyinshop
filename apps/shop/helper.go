package main

import (
	"context"
	"github.com/123508/douyinshop/pkg/myredis"
	"github.com/123508/douyinshop/pkg/util"
	"github.com/redis/go-redis/v9"
)

const (
	serviceName = "shop"
	product     = "product"
)

var Rds = connectWithRedis()

func connectWithRedis() *redis.Client {
	rds, err := myredis.InitRedis()
	if err != nil {
		util.LogError("打开Redis连接失败", "connectWithRedis", "", err)
	}
	return rds
}

func DelShopProductCache(ctx context.Context, rds *redis.Client, shopId uint64) {
	// 构建 pattern
	pattern := util.TakeKey(product, shopId, "*")
	var cursor uint64 = 0
	var batchSize int64 = 100 // 可调整

	for {
		keys, nextCursor, err := rds.Scan(ctx, cursor, pattern, batchSize).Result()
		if err != nil {
			util.LogError("redis扫描错误", "DelShopProductCache", "", err)
			return
		}
		if len(keys) > 0 {
			pipe := rds.Pipeline()
			for _, key := range keys {
				pipe.Del(ctx, key)
			}
			if _, err := pipe.Exec(ctx); err != nil {
				util.LogError("redis pipeline删除错误", "DelShopProductCache", "", err)
				return
			}
		}
		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}
}
