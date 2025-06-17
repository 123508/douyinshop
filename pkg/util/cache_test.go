// cache_test.go
package util

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/123508/douyinshop/pkg/models"
	"github.com/123508/douyinshop/pkg/myredis"
)

func TestSimpleCache(t *testing.T) {
	rds, err := myredis.InitRedis()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	com := SimpleCacheComponent[int64, string]{
		Rds:       rds,
		Ctx:       ctx,
		Key:       "Key A",
		FuncName:  "TestSimpleCache",
		Marshal:   json.Marshal,
		Unmarshal: json.Unmarshal,
		QueryExec: func() (string, error) {
			fmt.Println("缓存未命中，从数据库中查询")
			return "testdata A", nil
		},
		Expires: time.Minute * 5,
	}

	for i := 0; i < 10; i++ {
		time.Sleep(time.Second)
		data, err := com.QueryWithCache(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(data)
		if data != "testdata A" {
			t.Errorf("expected 'testdata A', got %v", data)
		}
	}
}

func TestListCache(t *testing.T) {
	rds, err := myredis.InitRedis()
	if err != nil {
		t.Fatal(err)
	}
	data := []models.User{
		{
			ID:     1,
			Name:   "User1",
			Email:  "user1@test.com",
			Phone:  "12300000000",
			Gender: 1,
			Avatar: "头像1",
			Status: 0,
		},
		{
			ID:     2,
			Name:   "User2",
			Email:  "user2@test.com",
			Phone:  "13300000000",
			Gender: 0,
			Avatar: "头像2",
			Status: 0,
		},
		{
			ID:     3,
			Name:   "User3",
			Email:  "user3@test.com",
			Phone:  "14300000000",
			Gender: 0,
			Avatar: "头像3",
			Status: 0,
		},
	}
	ctx := context.Background()
	com := ListCacheComponent[uint64, models.User]{
		Rds:             rds,
		Ctx:             ctx,
		IdListKey:       "IdListKey",
		DetailKeyPrefix: "DetailKeyPrefix:",
		FuncName:        "TestListCache",
		Marshal:         json.Marshal,
		Unmarshal:       json.Unmarshal,
		FullQueryExec: func() ([]models.User, error) {
			fmt.Println("全量查询未命中，从数据库中查询")
			return data, nil
		},
		PartialQueryExec: func(ids []uint64) ([]models.User, error) {
			fmt.Println("部分查询未命中，从数据库中查询")
			return data[1:2], nil
		},
		Expires:     time.Minute * 5,
		MaxLostRate: 30,
	}
	for i := 0; i < 10; i++ {
		items, err := com.QueryListWithCache(ctx)
		if err != nil {
			t.Fatal(err)
		}
		for j := 0; j < len(items); j++ {
			item := items[j]
			// 检查用户信息与data中是否对应
			if item.ID != data[j].ID || item.Name != data[j].Name || item.Email != data[j].Email || item.Phone != data[j].Phone || item.Gender != data[j].Gender || item.Avatar != data[j].Avatar || item.Status != data[j].Status {
				t.Errorf("Item %d does not match original data", j)
			}
			fmt.Printf("ID: %d, Name: %s, Email: %s, Phone: %s, Gender: %d, Avatar: %s, Status: %d\n", item.ID, item.Name, item.Email, item.Phone, item.Gender, item.Avatar, item.Status)
		}
	}
}
