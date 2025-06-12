package main

import (
	"context"
	"encoding/json"
	"github.com/123508/douyinshop/kitex_gen/address"
	"github.com/123508/douyinshop/pkg/db"
	"github.com/123508/douyinshop/pkg/errorno"
	"github.com/123508/douyinshop/pkg/models"
	"github.com/123508/douyinshop/pkg/myredis"
	"github.com/123508/douyinshop/pkg/util"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"math/rand"
	"strconv"
	"time"
)

//缓存策略采取先写入数据库后写入缓存,这种策略存在的问题是缓存过时失效
//同时采用随机过期策略防止雪崩
//增删改数据只将结果写入缓存,除此之外不涉及任何缓存有关的操作
//查找需要先检查缓存,如果有就返回,没有查完数据库后写入缓存再返回

// AddressServiceImpl implements the last service interface defined in the IDL.
type AddressServiceImpl struct{}

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

var rds = connectWithRedis()

func connectWithRedis() *redis.Client {
	rds, err := myredis.InitRedis()
	if err != nil {
		util.LogError("打开Redis连接失败", "connectWithRedis", "", err)
	}
	return rds
}

var InvalidAddressIdError = &errorno.BasicMessageError{Code: 400, Message: "地址ID无效"}

var ForbiddenDeleteError = &errorno.BasicMessageError{Code: 401, Message: "地址不存在或无权限删除"}

var DeleteAddrError = &errorno.BasicMessageError{Code: 500, Message: "地址删除失败,请联系管理员"}

var ForbiddenAskError = &errorno.BasicMessageError{Code: 401, Message: "无权访问该地址"}

var FailUpdateError = &errorno.BasicMessageError{Code: 500, Message: "更新地址失败,请联系管理员"}

var GetDefaultError = &errorno.BasicMessageError{Code: 404, Message: "获取默认地址失败"}

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

// 获取默认地址
func (s *AddressServiceImpl) getDefaultAddress(ctx context.Context, UserId uint64) uint64 {

	//查询失败,则重新查询数据库
	var item models.AddressBook
	DB.Model(&models.AddressBook{}).Where("user_id = ? and is_default= ?", UserId, true).First(&item)

	//没有用户默认地址,返回0
	return uint64(item.ID)
}

// AddAddress implements the AddressServiceImpl interface.
// 增加地址接口
// 如果一开始没有地址,就将新地址设置为默认地址,否则不为默认地址
func (s *AddressServiceImpl) AddAddress(ctx context.Context, req *address.AddAddressReq) (resp *address.AddAddressResp, err error) {

	//设置默认地址必须存在的逻辑,如果需要可以打开
	//listResp, err := s.GetAddressList(ctx, &address.GetAddressListReq{UserId: req.UserId})
	//
	//if err != nil {
	//	return nil, err
	//}
	//address1.IsDefault = len(listResp.Address) == 0

	//删除对应缓存
	defer util.CleanCache(rds, ctx, util.TakeKey(serviceName, "default", req.UserId))

	addr := &models.AddressBook{}
	addr = tranAddressToAddressBook(req.Address)

	err = DB.Transaction(func(tx *gorm.DB) error {
		//如果设置当前地址为默认地址
		if req.Address.IsDefault {
			//查询默认地址是否存在
			defaultAddressId := s.getDefaultAddress(ctx, req.UserId)
			//如果存在就先将默认地址取消
			if defaultAddressId != 0 {
				if err := DB.Model(&models.AddressBook{}).Where("id = ?", defaultAddressId).Update("is_default", false).Error; err != nil {
					util.LogError("变更默认地址失败", "AddAddress", "更新默认地址", err)
					return err
				}
			}

			addr.IsDefault = true
		}
		//创建新地址
		if err := DB.Create(&addr).Update("user_id", req.UserId).Error; err != nil {
			util.LogError("创建新地址失败", "AddAddress", "新建用户地址", err)
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &address.AddAddressResp{AddrId: uint64(addr.ID)}, nil
}

// DeleteAddress implements the AddressServiceImpl interface.
// 删除地址接口
func (s *AddressServiceImpl) DeleteAddress(ctx context.Context, req *address.DeleteAddressReq) (resp *address.DeleteAddressResp, err error) {

	//删除对应缓存
	defer func() {
		util.CleanCache(rds, ctx, util.TakeKey(serviceName, info, req.AddrId))
		util.CleanCache(rds, ctx, util.TakeKey(serviceName, "default", req.UserId))
		DelUserOrderCache(ctx, rds, req.UserId)
	}()

	// 验证地址ID是否有效
	if req.AddrId == 0 {
		return nil, InvalidAddressIdError
	}

	//默认地址必须存在的逻辑,如果需要可以打开
	//if s.getDefaultAddress(req.UserId) == req.AddrId {
	//	return &address.DeleteAddressResp{Res: false}, &errors.BasicMessageError{Message: "无法删除默认地址,请修改默认地址后删除该地址"}
	//}

	// 开始事务
	tx := DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 检查地址是否属于用户
	var addr models.AddressBook
	if err = tx.Where("id = ? and user_id = ?", req.AddrId, req.UserId).First(&addr).Error; err != nil {
		tx.Rollback()
		return &address.DeleteAddressResp{Res: false}, ForbiddenAskError
	}

	// 执行删除操作
	if err = tx.Unscoped().Delete(&addr).Error; err != nil {
		tx.Rollback()
		return &address.DeleteAddressResp{Res: false}, DeleteAddrError
	}

	// 提交事务
	if err = tx.Commit().Error; err != nil {
		util.LogError("事务出错", "DeleteAddress", "提交事务", err)
		return &address.DeleteAddressResp{Res: false}, DeleteAddrError
	}

	// 返回成功响应
	return &address.DeleteAddressResp{Res: true}, nil
}

// UpdateAddress implements the AddressServiceImpl interface.
// 更新地址接口
func (s *AddressServiceImpl) UpdateAddress(ctx context.Context, req *address.UpdateAddressReq) (resp *address.UpdateAddressResp, err error) {

	//删除对应缓存
	defer util.CleanCache(rds, ctx, util.TakeKey(serviceName, info, req.AddrId))
	defer util.CleanCache(rds, ctx, util.TakeKey(serviceName, "default", req.UserId))

	tx := DB.Model(&models.AddressBook{}).Where("id = ?", req.AddrId)

	var row models.AddressBook
	//如果更新地址不存在就直接返回错误(异常状态)
	if tx.Find(&row); row.ID == 0 {
		return &address.UpdateAddressResp{Res: false}, InvalidAddressIdError
	}

	// 检查地址是否属于用户
	var addr models.AddressBook
	if err = DB.Model(&models.AddressBook{}).Where("id = ? and user_id = ?", req.AddrId, req.UserId).First(&addr).Error; err != nil {
		return nil, ForbiddenDeleteError
	}

	err = DB.Transaction(func(tx *gorm.DB) error {
		tx = DB.Model(&models.AddressBook{}).Where("id = ?", req.AddrId)

		updates := make(map[string]interface{}, 10)

		//更新StreetAddress字段
		if req.Address.StreetAddress != "" {
			updates["stress_address"] = req.Address.StreetAddress
		}
		//更新Phone字段
		if req.Address.Phone != "" {
			updates["phone"] = req.Address.Phone
		}
		//更新.ZipCode字段
		if req.Address.ZipCode != "" {
			updates["zip_code"] = req.Address.ZipCode
		}
		//更新State字段
		if req.Address.State != "" {
			updates["state"] = req.Address.State
		}
		//更新City字段
		if req.Address.City != "" {
			updates["city"] = req.Address.City
		}
		//更新Consignee字段
		if req.Address.Consignee != "" {
			updates["consignee"] = req.Address.Consignee
		}
		//更新Country字段
		if req.Address.Country != "" {
			updates["country"] = req.Address.Country
		}
		//更新Label字段
		if req.Address.Label != "" {
			updates["label"] = req.Address.Label
		}
		//更新Gender字段
		if row.Gender != req.Address.Gender {
			updates["gender"] = req.Address.Gender
		}

		//更新IsDefault字段
		if req.Address.IsDefault {
			if _, err := s.SetDefaultAddress(context.Background(), &address.SetDefaultAddressReq{UserId: req.UserId, AddrId: req.AddrId}); err != nil {
				return err
			}
		} else {
			//默认地址可以不存在的逻辑,如果需要可以关闭
			updates["is_default"] = false
		}
		if err := tx.Updates(updates).Error; err != nil {
			return err
		}

		//事务提交
		return nil
	})

	if err != nil {
		util.LogError("事务出错", "UpdateAddress", "提交事务", err)
		return &address.UpdateAddressResp{Res: false}, FailUpdateError
	}

	return &address.UpdateAddressResp{Res: true}, nil
}

// SetDefaultAddress implements the AddressServiceImpl interface.
// 设置默认地址
func (s *AddressServiceImpl) SetDefaultAddress(ctx context.Context, req *address.SetDefaultAddressReq) (resp *address.SetDefaultAddressResp, err error) {

	//删除缓存
	defer rds.Del(ctx, util.TakeKey(serviceName, "default", req.UserId))

	defaultAddressId := s.getDefaultAddress(ctx, req.UserId)

	// 检查地址是否属于用户
	var addr models.AddressBook
	if err = DB.Where("id = ? and user_id = ?", req.AddrId, req.UserId).First(&addr).Error; err != nil {
		return nil, ForbiddenAskError
	}

	//如果要设置的地址就是默认地址,直接返回(异常情况)
	if defaultAddressId == req.AddrId {
		return &address.SetDefaultAddressResp{Res: true}, nil
	}

	err = DB.Transaction(func(tx *gorm.DB) error {

		if err := DB.Model(&models.AddressBook{}).Where("id = ?", defaultAddressId).Update("is_default", false).Error; err != nil {
			return err
		}

		if err := DB.Model(&models.AddressBook{}).Where("id = ?", req.AddrId).Update("is_default", true).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		util.LogError("事务出错", "SetDefaultAddress", "提交事务", err)
		return &address.SetDefaultAddressResp{Res: false}, err
	}

	return &address.SetDefaultAddressResp{Res: true}, nil
}

// GetAddressList implements the AddressServiceImpl interface.
// 获取地址列表接口
func (s *AddressServiceImpl) GetAddressList(ctx context.Context, req *address.GetAddressListReq) (resp *address.GetAddressListResp, err error) {

	////通用组件缓存处理
	//component := util.ListCacheComponent[uint, models.AddressBook]{
	//	Rds:             rds,
	//	Ctx:             ctx,
	//	IdListKey:       util.TakeKey(serviceName, req.UserId,util.Md5Hash(util.TakeKey("占位符,后续有扩充参数请修改"))),
	//	DetailKeyPrefix: util.TakeKey(serviceName, info),
	//	Marshal:         json.Marshal,
	//	Unmarshal:       json.Unmarshal,
	//	FullQueryExec: func() ([]models.AddressBook, error) {
	//		res := make([]models.AddressBook, 0)
	//		// 执行查询并处理错误
	//		if err := DB.Model(&models.AddressBook{}).Where("user_id = ?", req.UserId).Find(&res).Error; err != nil {
	//			// 处理错误，例如返回错误或记录日志
	//			return nil, err
	//		}
	//		return res, nil
	//	},
	//
	//	PartialQueryExec: func(fail []uint64) ([]models.AddressBook, error) {
	//		missedAddress := make([]models.AddressBook, 0)
	//		if err := DB.Where("id IN ?", fail).Find(&missedAddress).Error; err != nil {
	//			util.LogError("查询Address错误", "GetAddressList", "", err)
	//			return nil, InvalidAddressIdError
	//		}
	//		return missedAddress, nil
	//	},
	//	Expires:     time.Duration(rand.Intn(3)+3) * time.Minute,
	//	MaxLostRate: 30,
	//}
	//
	//res, err := component.QueryListWithCache()
	//
	//if err != nil {
	//	return nil, err
	//}

	//特化对应缓存组件处理
	res := make([]models.AddressBook, 0)

	key := util.TakeKey(serviceName, req.UserId, util.Md5Hash(util.TakeKey("占位符,后续有扩充参数请修改")))

	jsonData, _ := rds.Get(ctx, key).Result()

	list := make([]uint64, 0)
	fail := make([]uint64, 0)
	addressMap := make(map[uint64]models.AddressBook)

	ok := json.Unmarshal([]byte(jsonData), &list)

	//按照list查询详情缓存
	for _, v := range list {
		t := models.AddressBook{}
		key := util.TakeKey(serviceName, info, v)
		jsonData, err := rds.Get(ctx, key).Result()
		if err != nil || json.Unmarshal([]byte(jsonData), &t) != nil {
			fail = append(fail, v)
			continue
		}
		addressMap[v] = t
	}

	//计算缓存失效比率

	rate := 100

	if len(list) != 0 {
		rate = len(fail) * 100 / len(list)
	}

	if ok != nil || len(res) == 0 || rate > 30 {
		// 执行查询并处理错误
		if err := DB.Model(&models.AddressBook{}).Where("user_id = ?", req.UserId).Find(&res).Error; err != nil {
			// 处理错误，例如返回错误或记录日志
			return nil, err
		}

		//构建地址ID数组
		idList := make([]uint64, 0, len(res))
		for _, v := range res {
			idList = append(idList, v.ID)
		}

		//序列化并存储ID列表缓存
		jsonData, err := json.Marshal(&idList)

		if err != nil {
			util.LogError("序列化地址数组错误", "GetAddressList", "请求hash为"+key, err)
		} else {
			// 设置随机过期时间，30~45 分钟
			if setErr := rds.Set(ctx, key, jsonData, time.Duration(rand.Intn(3)+3)*time.Minute).Err(); setErr != nil {
				util.LogError("存入address缓存失败", "GetAddressList", "", setErr)
			}
		}

		//地址详情分级存储
		for _, v := range res {
			key := util.TakeKey(serviceName, info, v.ID)
			jsonData, err := json.Marshal(&v)
			if err != nil {
				util.LogError("序列化地址错误", "GetAddressList", "address的id为"+strconv.Itoa(int(v.ID)), err)
			} else {
				if err = rds.Set(ctx, key, jsonData, time.Duration(rand.Intn(3)+3)*time.Minute).Err(); err != nil {
					util.LogError("缓存地址失败", "GetAddressList", "address的id为"+strconv.Itoa(int(v.ID)), err)
				}
			}
		}

	} else {

		//缓存失效比例较低,逐条查询并放入缓存
		if len(fail) > 0 {

			missedAddress := make([]models.AddressBook, 0)
			if err := DB.Where("id IN ?", fail).Find(&missedAddress).Error; err != nil {
				util.LogError("查询Address错误", "GetAddressList", "", err)
				return nil, InvalidAddressIdError
			}
			for _, o := range missedAddress {
				addressMap[o.ID] = o
				//放入缓存
				key := util.TakeKey(serviceName, info, o.ID)
				jsonData, err := json.Marshal(&o)
				if err != nil {
					util.LogError("序列化地址错误", "GetAddressList", "地址id为"+strconv.Itoa(int(o.ID)), err)
				} else {
					if err = rds.Set(ctx, key, jsonData, time.Duration(rand.Intn(3)+3)*time.Minute).Err(); err != nil {
						util.LogError("缓存地址失败", "GetAddressList", "地址id为"+strconv.Itoa(int(o.ID)), err)
					}
				}
			}
		}

		res = make([]models.AddressBook, 0)

		for _, v := range list {
			if v, ok := addressMap[v]; ok {
				res = append(res, v)
			}
		}

		log.WithFields(log.Fields{
			"方法名": "GetAddressList",
		}).Info("查询缓存成功")
	}

	//提前给定切片容量,优化性能
	result := make([]*address.AddressItem, 0, len(res))

	for _, k := range res {
		result = append(result, &address.AddressItem{AddrId: k.ID, Address: tranAddressBookToAddress(&k)})
	}

	_ = res

	return &address.GetAddressListResp{Address: result}, nil
}

// GetAddressInfo implements the AddressServiceImpl interface.
// 获取指定地址信息
func (s *AddressServiceImpl) GetAddressInfo(ctx context.Context, req *address.GetAddressInfoReq) (resp *address.GetAddressInfoResp, err error) {

	var addr models.AddressBook

	//查询缓存
	key := util.TakeKey(serviceName, info, req.AddrId)

	jsonData, _ := rds.Get(ctx, key).Result()

	ok := json.Unmarshal([]byte(jsonData), &addr)

	//缓存没有查到,直接查找数据库
	if ok != nil {
		if err = DB.Where("id = ? and user_id = ?", req.AddrId, req.UserId).First(&addr).Error; err != nil {
			return nil, ForbiddenAskError
		}

		jsonData, err := json.Marshal(&addr)

		if err != nil {
			util.LogError("序列化地址错误", "GetAddressInfo", "地址id为"+strconv.Itoa(int(addr.ID)), err)
		} else {
			if err = rds.Set(ctx, key, jsonData, time.Duration(rand.Intn(3)+3)*time.Minute).Err(); err != nil {
				util.LogError("缓存地址失败", "GetAddressInfo", "地址id为"+strconv.Itoa(int(addr.ID)), err)
			}
		}

	} else {
		log.WithFields(log.Fields{
			"方法名": "GetAddressInfo",
		}).Info("查询缓存成功")
	}

	if addr.ID == 0 {
		return nil, InvalidAddressIdError
	}

	return &address.GetAddressInfoResp{Addr: &address.Address{
		StreetAddress: addr.StressAddress,
		City:          addr.City,
		State:         addr.State,
		Country:       addr.Country,
		ZipCode:       addr.ZipCode,
		Consignee:     addr.Consignee,
		Gender:        addr.Gender,
		Phone:         addr.Phone,
		Label:         addr.Label,
		IsDefault:     addr.IsDefault,
		AddressId:     addr.ID,
	}}, nil
}

// GetDefaultAddress implements the AddressServiceImpl interface.
// 获取默认地址
func (s *AddressServiceImpl) GetDefaultAddress(ctx context.Context, req *address.GetDefaultAddressReq) (resp *address.GetDefaultAddressResp, err error) {

	////通用组件缓存处理
	//simple := util.SimpleCacheComponent[uint64, models.AddressBook]{
	//	Rds:       rds,
	//	Ctx:       ctx,
	//	Key:       util.TakeKey(serviceName, "default", req.UserId),
	//	Marshal:   json.Marshal,
	//	Unmarshal: json.Unmarshal,
	//	QueryExec: func() (models.AddressBook, error) {
	//		DefaultId := s.getDefaultAddress(ctx, req.UserId)
	//		var addr models.AddressBook
	//		if err = DB.Where(" id = ?", DefaultId).First(&addr).Error; err != nil {
	//			log.Println(err)
	//			return models.AddressBook{}, GetDefaultError
	//		}
	//		return addr, nil
	//	},
	//	Expires: time.Duration(rand.Intn(3)+3) * time.Minute,
	//}
	//
	//addr, err := simple.QueryWithCache()
	//
	//if err != nil {
	//	return nil, err
	//}

	//特化对应缓存组件处理,先查缓存
	key := util.TakeKey(serviceName, "default", req.UserId)

	result, _ := rds.Get(ctx, key).Result()

	var addr models.AddressBook

	ok := json.Unmarshal([]byte(result), &addr)

	//缓存查询失败
	if ok != nil {

		DefaultId := s.getDefaultAddress(ctx, req.UserId)

		if err = DB.Where(" id = ?", DefaultId).First(&addr).Error; err != nil {
			klog.Fatal(err)
			return nil, GetDefaultError
		}

		//重新缓存
		jsonData, _ := json.Marshal(&addr)

		rds.SetEx(ctx, key, string(jsonData), time.Duration(rand.Intn(15)+30)*time.Minute)
	} else {
		log.WithFields(log.Fields{
			"方法名": "GetDefaultAddress",
		}).Info("查询缓存成功")

	}

	return &address.GetDefaultAddressResp{
		Addr: &address.Address{
			StreetAddress: addr.StressAddress,
			City:          addr.City,
			State:         addr.State,
			Country:       addr.Country,
			ZipCode:       addr.ZipCode,
			Consignee:     addr.Consignee,
			Gender:        addr.Gender,
			Phone:         addr.Phone,
			Label:         addr.Label,
			IsDefault:     addr.IsDefault,
			AddressId:     addr.ID,
		},
	}, nil
}
