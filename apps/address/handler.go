package main

import (
	"context"
	address "github.com/123508/douyinshop/kitex_gen/address"
	"github.com/123508/douyinshop/pkg/db"
	"github.com/123508/douyinshop/pkg/errors"
	"github.com/123508/douyinshop/pkg/models"
	"github.com/cloudwego/kitex/pkg/klog"
	"gorm.io/gorm"
	"log"
)

// AddressServiceImpl implements the last service interface defined in the IDL.
type AddressServiceImpl struct{}

var DB = open()

//注意地址类型有Address,AddressItem,AddressBook

func tranAddressToAddressBook(origin *address.Address) *models.AddressBook {
	address := &models.AddressBook{}
	address.StressAddress = origin.StreetAddress
	address.Phone = origin.Phone
	address.Gender = origin.Gender
	address.Consignee = origin.Consignee
	address.State = origin.State
	address.City = origin.City
	address.Country = origin.Country
	address.Label = origin.Label
	address.ZipCode = origin.ZipCode //共计九个字段
	return address
}

func tranAddressBookToAddress(origin *models.AddressBook) *address.Address {
	address := &address.Address{}
	address.StreetAddress = origin.StressAddress
	address.Phone = origin.Phone
	address.Gender = origin.Gender
	address.Consignee = origin.Consignee
	address.State = origin.State
	address.City = origin.City
	address.Country = origin.Country
	address.Label = origin.Label
	address.ZipCode = origin.ZipCode
	address.IsDefault = origin.IsDefault //共计十个字段
	return address
}

// 获取默认地址
func (s *AddressServiceImpl) getDefaultAddress(UserId uint32) uint64 {
	var item models.AddressBook
	DB.Model(&models.AddressBook{}).Where("user_id = ? and is_default= ?", UserId, true).First(&item)
	//没有用户默认地址,返回0
	return uint64(item.ID)
}

func open() *gorm.DB {
	DB, err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	return DB
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

	address1 := &models.AddressBook{}
	address1 = tranAddressToAddressBook(req.Address)

	err = DB.Transaction(func(tx *gorm.DB) error {
		//如果设置当前地址为默认地址
		if req.Address.IsDefault {
			//查询默认地址是否存在
			defaultAddressId := s.getDefaultAddress(req.UserId)
			//如果存在就先将默认地址取消
			if defaultAddressId != 0 {
				if err := DB.Model(&models.AddressBook{}).Where("id = ?", defaultAddressId).Update("is_default", false).Error; err != nil {
					return err
				}
			}
			address1.IsDefault = true
		}
		//创建新地址
		if err := DB.Create(&address1).Update("user_id", req.UserId).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &address.AddAddressResp{AddrId: uint64(address1.ID)}, nil
}

// GetAddressList implements the AddressServiceImpl interface.
// 获取地址列表接口
func (s *AddressServiceImpl) GetAddressList(ctx context.Context, req *address.GetAddressListReq) (resp *address.GetAddressListResp, err error) {

	var res []models.AddressBook

	// 执行查询并处理错误
	if err := DB.Model(&models.AddressBook{}).Where("user_id = ?", req.UserId).Find(&res).Error; err != nil {
		// 处理错误，例如返回错误或记录日志
		return nil, err
	}

	//提前给定切片容量,优化性能
	result := make([]*address.AddressItem, 0, len(res))

	for _, k := range res {

		result = append(result, &address.AddressItem{AddrId: uint64(k.ID), Address: tranAddressBookToAddress(&k)})

	}

	_ = res

	return &address.GetAddressListResp{Address: result}, nil
}

// DeleteAddress implements the AddressServiceImpl interface.
// 删除地址接口
func (s *AddressServiceImpl) DeleteAddress(ctx context.Context, req *address.DeleteAddressReq) (resp *address.DeleteAddressResp, err error) {

	// 验证地址ID是否有效
	if req.AddrId == 0 {
		return nil, &errors.BasicMessageError{Message: "地址ID无效"}
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
	var address1 models.AddressBook
	if err = tx.Where("id = ? and user_id = ?", req.AddrId, req.UserId).First(&address1).Error; err != nil {
		tx.Rollback()
		return &address.DeleteAddressResp{Res: false}, &errors.BasicMessageError{Message: "地址不存在或无权限删除"}
	}

	// 执行删除操作
	if err = tx.Unscoped().Delete(&address1).Error; err != nil {
		tx.Rollback()
		return &address.DeleteAddressResp{Res: false}, &errors.BasicMessageError{Message: "地址删除失败,请联系管理员"}
	}

	// 提交事务
	if err = tx.Commit().Error; err != nil {
		return &address.DeleteAddressResp{Res: false}, &errors.BasicMessageError{Message: "地址删除失败,请联系管理员"}
	}

	// 返回成功响应
	return &address.DeleteAddressResp{Res: true}, nil
}

// UpdateAddress implements the AddressServiceImpl interface.
// 更新地址接口
func (s *AddressServiceImpl) UpdateAddress(ctx context.Context, req *address.UpdateAddressReq) (resp *address.UpdateAddressResp, err error) {

	tx := DB.Model(&models.AddressBook{}).Where("id = ?", req.AddrId)

	var row models.AddressBook
	//如果更新地址不存在就直接返回错误(异常状态)
	if tx.Find(&row); row.ID == 0 {
		return &address.UpdateAddressResp{Res: false}, &errors.BasicMessageError{Message: "地址不存在,请联系管理员"}
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
		if req.Address.ZipCode != 0 {
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
		klog.Fatal(err)
		return &address.UpdateAddressResp{Res: false}, &errors.BasicMessageError{Message: "更新地址失败,请联系管理员"}
	}

	return &address.UpdateAddressResp{Res: true}, nil
}

// SetDefaultAddress implements the AddressServiceImpl interface.
// 设置默认地址
func (s *AddressServiceImpl) SetDefaultAddress(ctx context.Context, req *address.SetDefaultAddressReq) (resp *address.SetDefaultAddressResp, err error) {
	defaultAddressId := s.getDefaultAddress(req.UserId)

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
		return &address.SetDefaultAddressResp{Res: false}, err
	}

	return &address.SetDefaultAddressResp{Res: true}, nil
}
