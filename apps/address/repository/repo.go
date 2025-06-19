package repository

import (
	"context"
	"fmt"
	"github.com/123508/douyinshop/pkg/models"
	"github.com/123508/douyinshop/pkg/util"
	"gorm.io/gorm"
)

type AddressRepository interface {
	GetDB() *gorm.DB
	CreateAddress(ctx context.Context, address *models.AddressBook) (addrId uint64, err error)
	GetAddressById(ctx context.Context, UserId uint64, AddrId uint64) (address *models.AddressBook, err error)
	GetAddressList(ctx context.Context, page, pageSize int, userId uint64) (list []models.AddressBook, err error)
	UpdateAddress(ctx context.Context, address *models.AddressBook) (err error)
	DeleteAddress(ctx context.Context, AddrId uint64, UserId uint64) (err error)
	GetDefaultAddress(ctx context.Context, UserId uint64) (address *models.AddressBook, err error)
	RemoveDefaultAddress(ctx context.Context, UserId uint64) (err error)
	SetDefaultAddress(ctx context.Context, AddrId uint64, UserId uint64) (err error)
	AskAddress(ctx context.Context, UserId uint64, AddrId uint64) (is bool, address *models.AddressBook, err error)
}

type RepoImpl struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) AddressRepository {
	return &RepoImpl{
		DB: db,
	}
}

func (r *RepoImpl) GetDB() *gorm.DB {
	return r.DB
}

func (r *RepoImpl) CreateAddress(ctx context.Context, address *models.AddressBook) (addrId uint64, err error) {

	if err = r.DB.WithContext(ctx).Create(&address).Error; err != nil {
		util.LogError("创建新地址失败", "AddAddress", err)
		return 0, err
	}

	return address.ID, nil
}

func (r *RepoImpl) GetAddressById(ctx context.Context, UserId, AddrId uint64) (address *models.AddressBook, err error) {

	addr := models.AddressBook{}

	if err = r.DB.WithContext(ctx).Where("id = ? and user_id = ?", AddrId, UserId).First(&addr).Error; err != nil {
		util.LogError("获取地址失败", "GetAddressById", err)
		return nil, err
	}

	return &addr, nil
}

func (r *RepoImpl) GetAddressList(ctx context.Context, page, pageSize int, userId uint64) (list []models.AddressBook, err error) {

	offset := (page - 1) * pageSize

	list = make([]models.AddressBook, 0)

	if err = r.DB.
		Model(&models.AddressBook{}).
		WithContext(ctx).
		Where("user_id=?", userId).
		Offset(offset).
		Limit(pageSize).
		Find(&list).Error; err != nil {
		util.LogError("获取地址列表失败", "GetAddressList", err)
		return nil, err
	}

	return list, nil
}

func (r *RepoImpl) UpdateAddress(ctx context.Context, address *models.AddressBook) (err error) {

	updates := make(map[string]interface{}, 10)

	//更新StreetAddress字段
	if address.StressAddress != "" {
		updates["stress_address"] = address.StressAddress
	}
	//更新Phone字段
	if address.Phone != "" {
		updates["phone"] = address.Phone
	}
	//更新.ZipCode字段
	if address.ZipCode != "" {
		updates["zip_code"] = address.ZipCode
	}
	//更新State字段
	if address.State != "" {
		updates["state"] = address.State
	}
	//更新City字段
	if address.City != "" {
		updates["city"] = address.City
	}
	//更新Consignee字段
	if address.Consignee != "" {
		updates["consignee"] = address.Consignee
	}
	//更新Country字段
	if address.Country != "" {
		updates["country"] = address.Country
	}
	//更新Label字段
	if address.Label != "" {
		updates["label"] = address.Label
	}
	//更新Gender字段
	updates["gender"] = address.Gender

	//更新IsDefault字段
	updates["is_default"] = address.IsDefault

	//fmt.Println(address)

	if err := r.DB.
		Model(&models.AddressBook{}).
		WithContext(ctx).
		Where("id = ?", address.ID).
		Updates(updates).
		Error; err != nil {
		util.LogError("更新地址失败", "UpdateAddress", err)
		return err
	}

	return nil
}

func (r *RepoImpl) DeleteAddress(ctx context.Context, AddrId, UserId uint64) (err error) {
	if err = r.DB.
		Model(&models.AddressBook{}).
		WithContext(ctx).
		Where("id= ? and user_id = ?", AddrId, UserId).
		Unscoped().
		Delete(&models.AddressBook{ID: AddrId}).
		Error; err != nil {
		util.LogError("删除地址失败", "DeleteAddress", err)
		return err
	}
	return nil
}

func (r *RepoImpl) GetDefaultAddress(ctx context.Context, UserId uint64) (addressBook *models.AddressBook, err error) {

	var item models.AddressBook

	if err = r.DB.
		Model(&models.AddressBook{}).
		WithContext(ctx).
		Where("user_id = ? and is_default= ?", UserId, true).
		First(&item).
		Error; err != nil {
		util.LogError("获取默认地址失败", "GetDefaultAddress", err)
		return nil, err
	}

	//没有用户默认地址,返回0
	return &item, nil
}

func (r *RepoImpl) SetDefaultAddress(ctx context.Context, AddrId, UserId uint64) error {
	tx := r.DB.Begin().WithContext(ctx)

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// 1. 清除旧默认地址 (已包含用户ID)
	if err := tx.Model(&models.AddressBook{}).
		Where("user_id = ? AND is_default = true", UserId).
		Update("is_default", false).Error; err != nil {

		tx.Rollback()
		util.LogError("清除默认地址失败", "SetDefaultAddress", err)
		return fmt.Errorf("清除默认地址失败: %w", err)
	}

	// 2. 设置新默认地址 (关键校验)
	result := tx.Model(&models.AddressBook{}).
		Where("id = ? AND user_id = ?", AddrId, UserId).
		Update("is_default", true)

	if err := result.Error; err != nil {
		tx.Rollback()
		util.LogError("设置默认地址失败", "SetDefaultAddress", err)
		return fmt.Errorf("设置默认地址失败: %w", err)
	}

	// 显式检查是否更新成功
	if result.RowsAffected == 0 {
		tx.Rollback()
		util.LogError("地址不存在或不属于用户", "SetDefaultAddress", nil)
		return fmt.Errorf("地址不存在或不属于当前用户")
	}

	if err := tx.Commit().Error; err != nil {
		util.LogError("事务提交失败", "SetDefaultAddress", err)
		return fmt.Errorf("事务提交失败: %w", err)
	}

	return nil
}

func (r *RepoImpl) RemoveDefaultAddress(ctx context.Context, UserId uint64) (err error) {

	if err = r.DB.
		Model(&models.AddressBook{}).
		WithContext(ctx).
		Where("user_id = ? and is_default= ?", UserId, true).
		Update("is_default", false).
		Error; err != nil {
		util.LogError("移除默认地址失败", "RemoveDefaultAddress", err)
		return err
	}

	return nil
}

func (r *RepoImpl) AskAddress(ctx context.Context, UserId uint64, AddrId uint64) (is bool, address *models.AddressBook, err error) {
	var item models.AddressBook
	if err = r.DB.
		Model(&models.AddressBook{}).
		WithContext(ctx).
		Where("user_id = ? and id= ?", UserId, AddrId).
		First(&item).
		Error; err != nil {
		util.LogError("查询是否为用户的地址失败", "AskAddress", err)
		return false, nil, err
	}

	if item.ID == 0 {
		return false, nil, nil
	}

	return true, &item, nil
}
