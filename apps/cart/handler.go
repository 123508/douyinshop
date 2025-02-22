package main

import (
	"context"

	cart "github.com/123508/douyinshop/kitex_gen/cart"
)

// CartServiceImpl implements the last service interface defined in the IDL.
type CartServiceImpl struct{}

// AddItem implements the CartServiceImpl interface.
func (s *CartServiceImpl) AddItem(ctx context.Context, req *cart.AddItemReq) (resp *cart.AddItemResp, err error) {
	// TODO: Your code here...
	return
}

// GetCart implements the CartServiceImpl interface.
func (s *CartServiceImpl) GetCart(ctx context.Context, req *cart.GetCartReq) (resp *cart.GetCartResp, err error) {
	// TODO: Your code here...
	return
}

// EmptyCart implements the CartServiceImpl interface.
func (s *CartServiceImpl) EmptyCart(ctx context.Context, req *cart.EmptyCartReq) (resp *cart.EmptyCartResp, err error) {
	// TODO: Your code here...
	return
}

// DeleteItem implements the CartServiceImpl interface.
func (s *CartServiceImpl) DeleteItem(ctx context.Context, req *cart.DeleteItemReq) (resp *cart.EmptyCartResp, err error) {
	// TODO: Your code here...
	return
}
