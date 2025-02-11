package main

import (
	"context"
	shop "github.com/123508/douyinshop/kitex_gen/shop"
)

// ShopServiceImpl implements the last service interface defined in the IDL.
type ShopServiceImpl struct{}

// Register implements the ShopServiceImpl interface.
func (s *ShopServiceImpl) Register(ctx context.Context, req *shop.RegisterShopReq) (resp *shop.RegisterShopResp, err error) {
	// TODO: Your code here...
	return
}

// GetShopId implements the ShopServiceImpl interface.
func (s *ShopServiceImpl) GetShopId(ctx context.Context, req *shop.GetShopIdReq) (resp *shop.GetShopIdResp, err error) {
	// TODO: Your code here...
	return
}

// GetShopInfo implements the ShopServiceImpl interface.
func (s *ShopServiceImpl) GetShopInfo(ctx context.Context, req *shop.GetShopInfoReq) (resp *shop.GetShopInfoResp, err error) {
	// TODO: Your code here...
	return
}

// UpdateShopInfo implements the ShopServiceImpl interface.
func (s *ShopServiceImpl) UpdateShopInfo(ctx context.Context, req *shop.UpdateShopInfoReq) (resp *shop.UpdateShopInfoResp, err error) {
	// TODO: Your code here...
	return
}

// AddProduct implements the ShopServiceImpl interface.
func (s *ShopServiceImpl) AddProduct(ctx context.Context, req *shop.AddProductReq) (resp *shop.AddProductResp, err error) {
	// TODO: Your code here...
	return
}

// DeleteProduct implements the ShopServiceImpl interface.
func (s *ShopServiceImpl) DeleteProduct(ctx context.Context, req *shop.DeleteProductReq) (resp *shop.DeleteProductResp, err error) {
	// TODO: Your code here...
	return
}

// UpdateProduct implements the ShopServiceImpl interface.
func (s *ShopServiceImpl) UpdateProduct(ctx context.Context, req *shop.UpdateProductReq) (resp *shop.UpdateProductResp, err error) {
	// TODO: Your code here...
	return
}

// GetProductList implements the ShopServiceImpl interface.
func (s *ShopServiceImpl) GetProductList(ctx context.Context, req *shop.GetProductListReq) (resp *shop.GetProductListResp, err error) {
	// TODO: Your code here...
	return
}
