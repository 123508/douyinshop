package client

import (
	"context"
	"github.com/123508/douyinshop/kitex_gen/cart"
	"github.com/123508/douyinshop/kitex_gen/cart/cartservice"

	"github.com/123508/douyinshop/pkg/config"

	"time"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/retry"
	etcd "github.com/kitex-contrib/registry-etcd"
)

var cartClient cartservice.Client

func initCartRpc() {
	r, err := etcd.NewEtcdResolverWithAuth(config.Conf.EtcdConfig.Endpoints, config.Conf.EtcdConfig.Username, config.Conf.EtcdConfig.Password)
	if err != nil {
		panic(err)
	}

	c, err := cartservice.NewClient(
		config.Conf.AIConfig.ServiceName,                  // service name
		client.WithRPCTimeout(3*time.Second),              // rpc timeout
		client.WithConnectTimeout(50*time.Millisecond),    // conn timeout
		client.WithFailureRetry(retry.NewFailurePolicy()), // retry
		client.WithResolver(r),                            // resolver
	)
	if err != nil {
		panic(err)
	}
	cartClient = c
}

// AddItem 添加商品到购物车
func AddItem(ctx context.Context, req *cart.AddItemReq) (*cart.AddItemResp, error) {
	resp, err := cartClient.AddItem(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetCart获取购物车内容
func GetCart(ctx context.Context, req *cart.GetCartReq) (*cart.GetCartResp, error) {
	resp, err := cartClient.GetCart(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// EmptyCart 清空购物车
func EmptyCart(ctx context.Context, req *cart.EmptyCartReq) (*cart.EmptyCartResp, error) {
	resp, err := cartClient.EmptyCart(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// DeleteItem 删除购物车商品
func DeleteItem(ctx context.Context, req *cart.DeleteItemReq) (*cart.EmptyCartResp, error) {
	resp, err := cartClient.DeleteItem(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
