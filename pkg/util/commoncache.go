package util

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"strconv"
	"time"
)

type IntegerNumber interface {
	~int | ~uint | ~int64 | ~uint64 | ~int32 | ~uint32 | ~int16 | ~uint16 | ~int8 | ~uint8
}

type ItemType[E IntegerNumber] interface {
	GetID() E
}

type ListCacheComponent[Id IntegerNumber, Item ItemType[Id]] struct {
	Rds              *redis.Client
	Ctx              context.Context
	IdListKey        string
	DetailKeyPrefix  string
	funcName         string
	Marshal          func(v any) ([]byte, error)
	Unmarshal        func(data []byte, target any) error
	FullQueryExec    func() ([]Item, error)
	PartialQueryExec func([]Id) ([]Item, error)
	Expires          time.Duration
	MaxLostRate      int
	Sort             func([]Item) []Item
}

func (c *ListCacheComponent[Id, Item]) checkAndRepair() error {

	if c.Expires <= 0 {
		c.Expires = time.Duration(rand.Intn(3)+3) * time.Minute
	}

	if c.IdListKey == "" {
		return errors.New("IdListKey is required")
	}

	if c.DetailKeyPrefix == "" {
		return errors.New("DetailKeyPrefix is required")
	}

	if c.MaxLostRate > 100 {
		c.MaxLostRate = 30
	}

	if c.funcName == "" {
		c.funcName = "List"
	}

	if c.Rds == nil {
		return errors.New("rds is nil")
	}

	if c.Ctx == nil {
		c.Ctx = context.Background()
	}

	if c.Marshal == nil {
		c.Marshal = json.Marshal
	}

	if c.Unmarshal == nil {
		c.Unmarshal = json.Unmarshal
	}

	if c.FullQueryExec == nil {
		return errors.New("fullQueryExec is nil")
	}

	if c.PartialQueryExec == nil {
		return errors.New("partialQueryExec is nil")
	}

	if c.Sort == nil {
		c.Sort = func(items []Item) []Item {
			return items
		}
	}

	return nil
}

func (c *ListCacheComponent[Id, Item]) QueryListWithCache() ([]Item, error) {

	err := c.checkAndRepair()

	if err != nil {
		return nil, err
	}

	res := make([]Item, 0)

	key := c.IdListKey

	marshalData, err := c.Rds.Get(c.Ctx, key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		LogError("读取ID列表缓存失败", c.funcName, "查询hash为"+key, err)
		// 可选择直接回源
	}

	list := make([]Id, 0)
	fail := make([]Id, 0)
	itemMap := make(map[Id]Item)

	ok := c.Unmarshal([]byte(marshalData), &list)

	//按照list查询详情缓存
	for _, id := range list {
		var t Item
		key := TakeKey(c.DetailKeyPrefix, id)
		marshal, err := c.Rds.Get(c.Ctx, key).Result()
		if err != nil || c.Unmarshal([]byte(marshal), &t) != nil {
			fail = append(fail, id)
			continue
		}
		itemMap[id] = t
	}

	//计算缓存失效比率

	rate := 100

	if len(list) != 0 {
		rate = len(fail) * 100 / len(list)
	}

	if ok != nil || len(res) == 0 || rate > c.MaxLostRate {
		// 执行查询并处理错误
		result, err := c.FullQueryExec()

		if err != nil {
			return nil, err
		} else {
			res = result
		}

		//构建ID数组
		idList := make([]Id, 0, len(res))
		for _, v := range res {
			idList = append(idList, v.GetID())
		}

		//序列化并存储ID列表缓存
		marshal, err := c.Marshal(idList)

		if err != nil {
			LogError("序列化数组错误", c.funcName, "请求hash为"+key, err)
		} else {
			// 设置过期时间
			if setErr := c.Rds.Set(c.Ctx, key, marshal, c.Expires+time.Duration(rand.Intn(10))*time.Minute).Err(); setErr != nil {
				LogError("缓存失败", c.funcName, "", setErr)
			}
		}

		//详情内容分级存储
		for _, v := range res {
			key := TakeKey(c.DetailKeyPrefix, v.GetID())
			jsonData, err := c.Marshal(v)
			if err != nil {
				LogError("序列化错误", c.funcName, "id为"+strconv.Itoa(int(v.GetID())), err)
			} else {
				if err = c.Rds.Set(c.Ctx, key, jsonData, c.Expires+time.Duration(rand.Intn(10))*time.Minute).Err(); err != nil {
					LogError("缓存失败", c.funcName, "id为"+strconv.Itoa(int(v.GetID())), err)
				}
			}
		}

	} else {

		//缓存失效比例较低,逐条查询并放入缓存
		if len(fail) > 0 {
			var missedItem []Item

			result, err := c.PartialQueryExec(fail)

			if err != nil {
				return nil, err
			} else {
				missedItem = result
			}

			for _, item := range missedItem {
				itemMap[item.GetID()] = item
				//放入缓存
				key := TakeKey(c.DetailKeyPrefix, item.GetID())
				jsonData, err := c.Marshal(item)
				if err != nil {
					LogError("序列化错误", c.funcName, "id为"+strconv.Itoa(int(item.GetID())), err)
				} else {
					if err = c.Rds.Set(c.Ctx, key, jsonData, c.Expires+time.Duration(rand.Intn(10))*time.Second).Err(); err != nil {
						LogError("缓存失败", c.funcName, "id为"+strconv.Itoa(int(item.GetID())), err)
					}
				}
			}
		}

		res = make([]Item, 0)

		for _, v := range list {
			if v, ok := itemMap[v]; ok {
				res = append(res, v)
			}
		}

		log.WithFields(log.Fields{
			"方法名": c.funcName,
		}).Info("查询缓存成功")
	}

	return c.Sort(res), nil
}

type SimpleCacheComponent[Id IntegerNumber, E any] struct {
	Rds       *redis.Client
	Ctx       context.Context
	Key       string
	funcName  string
	Marshal   func(v any) ([]byte, error)
	Unmarshal func(data []byte, target any) error
	QueryExec func() (E, error)
	Expires   time.Duration
}

func (c *SimpleCacheComponent[Id, E]) checkAndRepair() error {
	if c.Expires <= 0 {
		c.Expires = time.Duration(rand.Intn(3)+3) * time.Minute
	}
	if c.Marshal == nil {
		c.Marshal = json.Marshal
	}
	if c.Unmarshal == nil {
		c.Unmarshal = json.Unmarshal
	}

	if c.QueryExec == nil {
		return errors.New("queryExec is nil")
	}

	if c.Key == "" {
		return errors.New("key is empty")
	}
	return nil
}

func (c *SimpleCacheComponent[Id, E]) QueryWithCache() (E, error) {

	var items E

	result, err := c.Rds.Get(c.Ctx, c.Key).Result()

	//查询缓存失败
	if err != nil || c.Unmarshal([]byte(result), &items) != nil {

		exec, err := c.QueryExec()

		if err != nil {
			return items, err
		} else {
			items = exec
		}

		marshal, err := c.Marshal(items)

		if err != nil {
			LogError("序列化数据失败", c.funcName, "序列化", err)
		} else {
			if setErr := c.Rds.Set(c.Ctx, c.Key, marshal, time.Duration(rand.Intn(3)+3)*time.Minute).Err(); setErr != nil {
				LogError("缓存数据失败", c.funcName, "加入缓存", setErr)
			}
		}

	} else {
		log.WithFields(log.Fields{
			"方法名": c.funcName,
		}).Info("查询缓存成功")
	}

	return items, nil
}
