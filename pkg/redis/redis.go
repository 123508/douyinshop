package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/123508/douyinshop/pkg/config"
)

func InitRedis() (*redis.Client, error) {
	RDB := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Conf.RedisConfig.Host, config.Conf.RedisConfig.Port),
		Password: config.Conf.RedisConfig.Password,
		DB:       0,
		PoolSize: config.Conf.RedisConfig.PoolSize,
	})

	_, err := RDB.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	return RDB, nil
}
