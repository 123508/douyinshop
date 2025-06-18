package service

import (
	"context"
	"encoding/json"
	"github.com/123508/douyinshop/apps/address/repository"
	"github.com/123508/douyinshop/pkg/models"
	"github.com/123508/douyinshop/pkg/util"
	"github.com/redis/go-redis/v9"
	"math/rand"
	"time"
)

type AddressService interface {
	GetRedis() *redis.Client
	AddAddress(ctx context.Context, addr *models.AddressBook) (addrId uint64, err error)
	DeleteAddress(ctx context.Context, AddrId, UserId uint64) error
	UpdateAddress(ctx context.Context, addr *models.AddressBook) error
	SetDefaultAddress(ctx context.Context, AddrId uint64, UserId uint64) error
	GetAddressList(ctx context.Context, UserId uint64) ([]models.AddressBook, error)
	GetAddressInfo(ctx context.Context, AddrId, UserId uint64) (*models.AddressBook, error)
	GetDefaultAddress(ctx context.Context, UserId uint64) (*models.AddressBook, error)
}

type ServiceImpl struct {
	addrRepo repository.AddressRepository
	Rds      *redis.Client
}

func NewService(addrRepo repository.AddressRepository, redis *redis.Client) AddressService {
	return &ServiceImpl{
		addrRepo: addrRepo,
		Rds:      redis,
	}
}

const (
	serviceName = "address"
	info        = "info"
)

//缓存策略采取先写入数据库后写入缓存,这种策略存在的问题是缓存过时失效
//同时采用随机过期策略防止雪崩
//增删改数据只将结果写入缓存,除此之外不涉及任何缓存有关的操作
//查找需要先检查缓存,如果有就返回,没有查完数据库后写入缓存再返回

func (s *ServiceImpl) GetRedis() *redis.Client {
	return s.Rds
}

func (s *ServiceImpl) AddAddress(ctx context.Context, addr *models.AddressBook) (addrId uint64, err error) {

	//删除对应缓存
	defer util.CleanCache(s.Rds, ctx, util.TakeKey(serviceName, "default", addr.UserId))
	DelAddressCache(ctx, s.Rds, addr.UserId)

	if addr.IsDefault {
		if err := s.addrRepo.RemoveDefaultAddress(ctx, addr.UserId); err != nil {
			return 0, err
		}
	}

	id, err := s.addrRepo.CreateAddress(ctx, addr)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *ServiceImpl) DeleteAddress(ctx context.Context, AddrId, UserId uint64) error {
	//删除对应缓存
	defer func() {
		util.CleanCache(s.Rds, ctx, util.TakeKey(serviceName, info, AddrId))
		util.CleanCache(s.Rds, ctx, util.TakeKey(serviceName, "default", AddrId))
		DelAddressCache(ctx, s.Rds, UserId)
	}()

	// 验证地址ID是否有效
	if AddrId == 0 || UserId == 0 {
		return InvalidAddressIdError
	}

	// 执行删除操作
	if err := s.addrRepo.DeleteAddress(ctx, AddrId, UserId); err != nil {
		return DeleteAddrError
	}

	return nil
}

func (s *ServiceImpl) UpdateAddress(ctx context.Context, addr *models.AddressBook) error {
	//删除对应缓存
	defer util.CleanCache(s.Rds, ctx, util.TakeKey(serviceName, info, addr.ID))
	defer util.CleanCache(s.Rds, ctx, util.TakeKey(serviceName, "default", addr.UserId))

	addressBook, err := s.addrRepo.GetAddressById(ctx, addr.UserId, addr.ID)

	if err != nil {
		return err
	}

	if addressBook.UserId != addr.UserId {
		return ForbiddenAskError
	}

	if addr.IsDefault {
		if err := s.addrRepo.RemoveDefaultAddress(ctx, addr.UserId); err != nil {
			return err
		}
	}

	return s.addrRepo.UpdateAddress(ctx, addr)
}

func (s *ServiceImpl) SetDefaultAddress(ctx context.Context, AddrId, UserId uint64) error {
	//删除缓存
	defer s.Rds.Del(ctx, util.TakeKey(serviceName, "default", UserId))

	// 检查地址是否属于用户
	is, _, err := s.addrRepo.AskAddress(ctx, UserId, AddrId)

	if !is || err != nil {
		return ForbiddenAskError
	}

	if err = s.addrRepo.SetDefaultAddress(ctx, AddrId, UserId); err != nil {
		return FailUpdateError
	}

	return nil
}

func (s *ServiceImpl) GetAddressList(ctx context.Context, UserId uint64) ([]models.AddressBook, error) {
	//通用组件缓存处理
	component := util.ListCacheComponent[uint64, models.AddressBook]{
		Rds:             s.Rds,
		Ctx:             ctx,
		IdListKey:       util.TakeKey(serviceName, UserId, util.Md5Hash(util.TakeKey("占位符,后续有扩充参数请修改"))),
		DetailKeyPrefix: util.TakeKey(serviceName, info),
		FuncName:        "GetAddressList",
		Marshal:         json.Marshal,
		Unmarshal:       json.Unmarshal,
		FullQueryExec: func() ([]models.AddressBook, error) {

			list, err := s.addrRepo.GetAddressList(ctx, 0, 100, UserId)

			if err != nil {
				return nil, err
			}

			return list, nil

		},
		Expires:     time.Duration(rand.Intn(3)+3) * time.Minute,
		MaxLostRate: 30,
		Sort:        nil,
		IdName:      "id",
		DB:          s.addrRepo.GetDB(),
	}

	res, err := component.QueryListWithCache()

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *ServiceImpl) GetAddressInfo(ctx context.Context, AddrId, UserId uint64) (*models.AddressBook, error) {

	simple := util.SimpleCacheComponent[uint64, models.AddressBook]{
		Rds:       s.Rds,
		Ctx:       ctx,
		Key:       util.TakeKey(serviceName, info, AddrId),
		FuncName:  "GetAddressInfo",
		Marshal:   json.Marshal,
		Unmarshal: json.Unmarshal,
		QueryExec: func() (models.AddressBook, error) {
			res, err := s.addrRepo.GetAddressById(ctx, UserId, AddrId)
			if err != nil {
				return models.AddressBook{}, err
			}
			return *res, nil
		},
		Expires: time.Duration(rand.Intn(10)+3) * time.Minute,
	}

	cache, err := simple.QueryWithCache()

	if err != nil {
		return nil, InvalidAddressIdError
	}

	return &cache, nil

}

func (s *ServiceImpl) GetDefaultAddress(ctx context.Context, UserId uint64) (*models.AddressBook, error) {
	//通用组件缓存处理
	simple := util.SimpleCacheComponent[uint64, models.AddressBook]{
		Rds:       s.Rds,
		Ctx:       ctx,
		Key:       util.TakeKey(serviceName, "default", UserId),
		Marshal:   json.Marshal,
		Unmarshal: json.Unmarshal,
		QueryExec: func() (models.AddressBook, error) {
			address, err := s.addrRepo.GetDefaultAddress(ctx, UserId)
			if err != nil {
				return models.AddressBook{}, err
			}
			return *address, nil
		},
		Expires: time.Duration(rand.Intn(13)+3) * time.Minute,
	}

	addr, err := simple.QueryWithCache()

	if err != nil {
		return nil, GetDefaultError
	}

	return &addr, nil
}

func DelAddressCache(ctx context.Context, rds *redis.Client, userId interface{}) {
	// 构建 pattern
	pattern := util.TakeKey(serviceName, userId, "*")
	var cursor uint64 = 0
	var batchSize int64 = 100 // 每次扫描的数量，可根据实际情况调整

	for {
		// 使用 SCAN 命令遍历所有匹配的 key
		keys, nextCursor, err := rds.Scan(ctx, cursor, pattern, batchSize).Result()
		if err != nil {
			util.LogError("redis扫描错误", "DelAddressCache", err)
		}
		if len(keys) > 0 {
			// pipeline 批量删除，提升性能
			pipe := rds.Pipeline()
			for _, key := range keys {
				pipe.Del(ctx, key)
			}
			if _, err := pipe.Exec(ctx); err != nil {
				util.LogError("redis pipeline删除错误", "DelAddressCache", err)
			}
		}
		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}
}
