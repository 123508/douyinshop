package client

import (
	"github.com/123508/douyinshop/kitex_gen/shop"
	"github.com/123508/douyinshop/kitex_gen/shop/shopservice"

	"github.com/123508/douyinshop/pkg/config"

	"context"
	"time"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/retry"
	etcd "github.com/kitex-contrib/registry-etcd"
)

var shopClient shopservice.Client

func init() {
	initShopRpc()
}

func initShopRpc() {
	r, err := etcd.NewEtcdResolverWithAuth(config.Conf.EtcdConfig.Endpoints, config.Conf.EtcdConfig.Username, config.Conf.EtcdConfig.Password)
	if err != nil {
		panic(err)
	}

	c, err := shopservice.NewClient(
		config.Conf.ShopConfig.ServiceName,                // service name
		client.WithRPCTimeout(3*time.Second),              // rpc timeout
		client.WithConnectTimeout(50*time.Millisecond),    // conn timeout
		client.WithFailureRetry(retry.NewFailurePolicy()), // retry
		client.WithResolver(r),                            // resolver
	)
	if err != nil {
		panic(err)
	}
	shopClient = c
}

func RegisterShop(ctx context.Context, req *shop.RegisterShopReq) (*shop.RegisterShopResp, error) {
	resp, err := shopClient.Register(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func GetShopId(ctx context.Context, req *shop.GetShopIdReq) (*shop.GetShopIdResp, error) {
	resp, err := shopClient.GetShopId(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func GetShopInfo(ctx context.Context, req *shop.GetShopInfoReq) (*shop.GetShopInfoResp, error) {
	resp, err := shopClient.GetShopInfo(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func UpdateShopInfo(ctx context.Context, req *shop.UpdateShopInfoReq) (*shop.UpdateShopInfoResp, error) {
	resp, err := shopClient.UpdateShopInfo(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func AddProduct(ctx context.Context, req *shop.AddProductReq) (*shop.AddProductResp, error) {
	resp, err := shopClient.AddProduct(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func DeleteProduct(ctx context.Context, req *shop.DeleteProductReq) (*shop.DeleteProductResp, error) {
	resp, err := shopClient.DeleteProduct(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func UpdateProduct(ctx context.Context, req *shop.UpdateProductReq) (*shop.UpdateProductResp, error) {
	resp, err := shopClient.UpdateProduct(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func GetProductList(ctx context.Context, req *shop.GetProductListReq) (*shop.GetProductListResp, error) {
	resp, err := shopClient.GetProductList(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
