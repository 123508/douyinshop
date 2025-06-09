package util

import (
	"context"
	"github.com/redis/go-redis/v9"
	"strconv"
	"strings"
	"time"
)

// TakeKey 构建key函数,默认为data[0]:data[1]:....:data[n],注意如果有无法识别的类型会在原本的位置填入 ???
func TakeKey(data ...any) string {
	builder := strings.Builder{}
	for i, v := range data {
		switch v.(type) {
		case string:
			builder.WriteString(v.(string))
		case int32:
			builder.WriteString(strconv.FormatInt(int64(v.(int32)), 10))
		case uint32:
			builder.WriteString(strconv.FormatInt(int64(v.(uint32)), 10))
		case int64:
			builder.WriteString(strconv.FormatInt(v.(int64), 10))
		case uint64:
			builder.WriteString(strconv.FormatUint(v.(uint64), 10))
		case int:
			builder.WriteString(strconv.Itoa(v.(int)))
		case uint:
			builder.WriteString(strconv.FormatInt(int64(v.(uint)), 10))
		case []byte:
			builder.WriteString(string(v.([]byte)))
		case time.Time:
			builder.WriteString(strconv.FormatInt(v.(time.Time).Unix(), 10))
		case byte:
			builder.WriteString(strconv.FormatInt(int64(v.(byte)), 10))
		}
		if i != len(data)-1 {
			builder.WriteString(":")
		}
	}

	return builder.String()
}

// CleanCache 清理缓存
func CleanCache(rds *redis.Client, ctx context.Context, key string) {
	err := rds.Del(ctx, key).Err()
	if err != nil {
		LogError("删除缓存错误", "cleanCache", "", err)
	}
}
