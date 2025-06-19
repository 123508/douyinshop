package client

import (
	"context"
	"github.com/123508/douyinshop/kitex_gen/address"
	"github.com/123508/douyinshop/kitex_gen/address/addressservice"

	"github.com/123508/douyinshop/pkg/config"

	"time"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/retry"
	etcd "github.com/kitex-contrib/registry-etcd"
)

var addressClient addressservice.Client

func initAddressRpc() {
	r, err := etcd.NewEtcdResolverWithAuth(config.Conf.EtcdConfig.Endpoints, config.Conf.EtcdConfig.Username, config.Conf.EtcdConfig.Password)
	if err != nil {
		panic(err)
	}

	c, err := addressservice.NewClient(
		config.Conf.AddressConfig.ServiceName,             // service name
		client.WithRPCTimeout(3*time.Second),              // rpc timeout
		client.WithConnectTimeout(50*time.Millisecond),    // conn timeout
		client.WithFailureRetry(retry.NewFailurePolicy()), // retry
		client.WithResolver(r),                            // resolver
	)
	if err != nil {
		panic(err)
	}
	addressClient = c
}

type AddressItem struct {
	ID            uint64 `json:"addr_id"`
	StreetAddress string `json:"street_address"`
	City          string `json:"city"`
	State         string `json:"state"`
	Country       string `json:"country"`
	ZipCode       string `json:"zip_code"`
	Consignee     string `json:"consignee"`
	Gender        uint32 `json:"gender"`
	Phone         string `json:"phone"`
	Label         string `json:"label"`
	IsDefault     bool   `json:"is_default"`
}

// AddAddress 添加地址
// item 地址信息
// userID 用户ID
// 返回地址ID
func AddAddress(ctx context.Context, item *AddressItem, userID uint64) (uint64, error) {
	req := &address.AddAddressReq{
		Address: &address.Address{
			StreetAddress: item.StreetAddress,
			City:          item.City,
			State:         item.State,
			Country:       item.Country,
			ZipCode:       item.ZipCode,
			Consignee:     item.Consignee,
			Gender:        item.Gender,
			Phone:         item.Phone,
			Label:         item.Label,
		},
		TargetUserId: userID,
	}
	resp, err := addressClient.AddAddress(ctx, req)
	if err != nil {
		return 0, err
	}
	return resp.AddrId, nil
}

// GetAddressList 获取用户地址列表
// userID 用户ID
// 返回地址列表
func GetAddressList(ctx context.Context, userID uint64) ([]AddressItem, error) {
	req := &address.GetAddressListReq{
		TargetUserId: userID,
	}
	resp, err := addressClient.GetAddressList(ctx, req)
	if err != nil {
		return nil, err
	}
	var addresses []AddressItem
	for _, a := range resp.Address {
		addresses = append(addresses, AddressItem{
			ID:            a.AddressId,
			StreetAddress: a.StreetAddress,
			City:          a.City,
			State:         a.State,
			Country:       a.Country,
			ZipCode:       a.ZipCode,
			Consignee:     a.Consignee,
			Gender:        a.Gender,
			Phone:         a.Phone,
			Label:         a.Label,
			IsDefault:     a.IsDefault,
		})
	}
	return addresses, nil
}

// DeleteAddress 删除地址
// addrID 地址ID
// userID 用户ID
// 返回是否删除成功
func DeleteAddress(ctx context.Context, addrID int, userID uint64) (bool, error) {
	req := &address.DeleteAddressReq{
		AddrId:       uint64(addrID),
		TargetUserId: userID,
	}
	resp, err := addressClient.DeleteAddress(ctx, req)
	if err != nil {
		return false, err
	}
	return resp.Res, nil
}

// UpdateAddress 更新地址
// item 地址信息
// userID 用户ID
// addrID 地址ID
// 返回是否更新成功
func UpdateAddress(ctx context.Context, item *AddressItem, userID uint64) (bool, error) {
	req := &address.UpdateAddressReq{
		Address: &address.Address{
			StreetAddress: item.StreetAddress,
			City:          item.City,
			State:         item.State,
			Country:       item.Country,
			ZipCode:       item.ZipCode,
			Consignee:     item.Consignee,
			Gender:        item.Gender,
			Phone:         item.Phone,
			Label:         item.Label,
			IsDefault:     item.IsDefault,
		},
		TargetUserId: userID,
		AddrId:       item.ID,
	}
	resp, err := addressClient.UpdateAddress(ctx, req)
	if err != nil {
		return false, err
	}
	return resp.Res, nil
}

// SetDefaultAddress 设置默认地址
// addrID 地址ID
// userID 用户ID
// 返回是否设置成功
func SetDefaultAddress(ctx context.Context, addrID int, userID uint64) (bool, error) {
	req := &address.SetDefaultAddressReq{
		AddrId:       uint64(addrID),
		TargetUserId: userID,
	}
	resp, err := addressClient.SetDefaultAddress(ctx, req)
	if err != nil {
		return false, err
	}
	return resp.Res, nil
}

// GetAddressInfo
// 获取用户制定地址的详细信息
func GetAddressInfo(ctx context.Context, addrId, userId uint64) (*AddressItem, error) {
	req := &address.GetAddressInfoReq{
		TargetUserId: userId,
		AddrId:       addrId,
	}
	resp, err := addressClient.GetAddressInfo(ctx, req)
	if err != nil {
		return nil, err
	}

	a := resp.Addr

	return &AddressItem{
		ID:            a.AddressId,
		StreetAddress: a.StreetAddress,
		City:          a.City,
		State:         a.State,
		Country:       a.Country,
		ZipCode:       a.ZipCode,
		Consignee:     a.Consignee,
		Gender:        a.Gender,
		Phone:         a.Phone,
		Label:         a.Label,
		IsDefault:     a.IsDefault,
	}, nil
}
