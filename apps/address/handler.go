package main

import (
	"context"
	address "github.com/123508/douyinshop/kitex_gen/address"
)

// AddressServiceImpl implements the last service interface defined in the IDL.
type AddressServiceImpl struct{}

// AddAddress implements the AddressServiceImpl interface.
func (s *AddressServiceImpl) AddAddress(ctx context.Context, req *address.AddAddressReq) (resp *address.AddAddressResp, err error) {
	// TODO: Your code here...
	return
}

// GetAddressList implements the AddressServiceImpl interface.
func (s *AddressServiceImpl) GetAddressList(ctx context.Context, req *address.GetAddressListReq) (resp *address.GetAddressListResp, err error) {
	// TODO: Your code here...
	return
}

// DeleteAddress implements the AddressServiceImpl interface.
func (s *AddressServiceImpl) DeleteAddress(ctx context.Context, req *address.DeleteAddressReq) (resp *address.DeleteAddressResp, err error) {
	// TODO: Your code here...
	return
}

// UpdateAddress implements the AddressServiceImpl interface.
func (s *AddressServiceImpl) UpdateAddress(ctx context.Context, req *address.UpdateAddressReq) (resp *address.UpdateAddressResp, err error) {
	// TODO: Your code here...
	return
}

// SetDefaultAddress implements the AddressServiceImpl interface.
func (s *AddressServiceImpl) SetDefaultAddress(ctx context.Context, req *address.SetDefaultAddressReq) (resp *address.SetDefaultAddressResp, err error) {
	// TODO: Your code here...
	return
}
