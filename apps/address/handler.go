package main

import (
	"context"
	"errors"
	"github.com/123508/douyinshop/apps/address/repository"
	"github.com/123508/douyinshop/apps/address/service"
	"github.com/123508/douyinshop/kitex_gen/address"
	"github.com/123508/douyinshop/pkg/models"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// AddressServiceImpl implements the last service interface defined in the IDL.
type AddressServiceImpl struct {
	AddressService service.AddressService
}

func NewAddressServiceImpl(db *gorm.DB, redis *redis.Client) *AddressServiceImpl {
	return &AddressServiceImpl{
		AddressService: service.NewService(repository.NewRepository(db), redis),
	}
}

// AddAddress implements the AddressServiceImpl interface.
// 增加地址接口
// 如果一开始没有地址,就将新地址设置为默认地址,否则不为默认地址
func (s *AddressServiceImpl) AddAddress(ctx context.Context, req *address.AddAddressReq) (resp *address.AddAddressResp, err error) {

	addressBook := AddressToAddressBook(req.Address)

	addressBook.UserId = req.TargetUserId

	id, err := s.AddressService.AddAddress(ctx, addressBook, req.RequestUserId)

	if err != nil {
		return nil, err
	}

	return &address.AddAddressResp{AddrId: id}, nil
}

// DeleteAddress implements the AddressServiceImpl interface.
// 删除地址接口
func (s *AddressServiceImpl) DeleteAddress(ctx context.Context, req *address.DeleteAddressReq) (resp *address.DeleteAddressResp, err error) {

	if err = s.AddressService.DeleteAddress(ctx, req.AddrId, req.TargetUserId, req.RequestUserId); err != nil {
		return nil, err
	}

	return &address.DeleteAddressResp{Res: true}, nil
}

// UpdateAddress implements the AddressServiceImpl interface.
// 更新地址接口
func (s *AddressServiceImpl) UpdateAddress(ctx context.Context, req *address.UpdateAddressReq) (resp *address.UpdateAddressResp, err error) {

	addressBook := AddressToAddressBook(req.Address)

	addressBook.UserId = req.TargetUserId

	addressBook.ID = req.AddrId

	if err = s.AddressService.UpdateAddress(ctx, addressBook, req.RequestUserId); err != nil {
		return nil, err
	}

	return &address.UpdateAddressResp{Res: true}, nil
}

// SetDefaultAddress implements the AddressServiceImpl interface.
// 设置默认地址
func (s *AddressServiceImpl) SetDefaultAddress(ctx context.Context, req *address.SetDefaultAddressReq) (resp *address.SetDefaultAddressResp, err error) {

	if err = s.AddressService.SetDefaultAddress(ctx, req.AddrId, req.TargetUserId, req.RequestUserId); err != nil {
		return nil, err
	}

	return &address.SetDefaultAddressResp{Res: true}, nil
}

// GetAddressList implements the AddressServiceImpl interface.
// 获取地址列表接口
func (s *AddressServiceImpl) GetAddressList(ctx context.Context, req *address.GetAddressListReq) (resp *address.GetAddressListResp, err error) {

	res, err := s.AddressService.GetAddressList(ctx, req.TargetUserId, req.RequestUserId)

	if err != nil {
		return nil, err
	}

	//提前给定切片容量,优化性能
	result := make([]*address.Address, 0, len(res))

	for _, k := range res {
		result = append(result, AddressBookToAddress(&k))
	}

	_ = res

	return &address.GetAddressListResp{Address: result}, nil
}

// GetAddressInfo implements the AddressServiceImpl interface.
// 获取指定地址信息
func (s *AddressServiceImpl) GetAddressInfo(ctx context.Context, req *address.GetAddressInfoReq) (resp *address.GetAddressInfoResp, err error) {

	addr, err := s.AddressService.GetAddressInfo(ctx, req.AddrId, req.TargetUserId, req.RequestUserId)

	if addr == nil || addr.ID == 0 {
		return nil, err
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

	addr, err := s.AddressService.GetDefaultAddress(ctx, req.TargetUserId, req.RequestUserId)

	if err != nil {
		return nil, err
	}

	if addr == nil {
		return nil, errors.New("没有默认地址")
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

// AddressToAddressBook 两个地址转换函数,注意地址类型有Address,AddressBook
func AddressToAddressBook(origin *address.Address) *models.AddressBook {
	addr := &models.AddressBook{}
	addr.ID = origin.AddressId
	addr.StressAddress = origin.StreetAddress
	addr.Phone = origin.Phone
	addr.Gender = origin.Gender
	addr.Consignee = origin.Consignee
	addr.State = origin.State
	addr.City = origin.City
	addr.Country = origin.Country
	addr.Label = origin.Label
	addr.ZipCode = origin.ZipCode
	addr.IsDefault = origin.IsDefault //共计一个字段
	return addr
}

func AddressBookToAddress(origin *models.AddressBook) *address.Address {
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
	addr.IsDefault = origin.IsDefault
	addr.AddressId = origin.ID //共计十一个字段
	return addr
}
