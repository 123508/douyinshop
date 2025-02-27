package client

import (
	"github.com/123508/douyinshop/kitex_gen/ai/aiservice"
	"github.com/123508/douyinshop/kitex_gen/ai"
	"github.com/123508/douyinshop/pkg/config"
	"context"
	"time"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/retry"
	etcd "github.com/kitex-contrib/registry-etcd"
)

var aiClient aiservice.Client

func initAiRpc() {
	r, err := etcd.NewEtcdResolverWithAuth(config.Conf.EtcdConfig.Endpoints, config.Conf.EtcdConfig.Username, config.Conf.EtcdConfig.Password)
	if err != nil {
		panic(err)
	}

	c, err := aiservice.NewClient(
		config.Conf.AIConfig.ServiceName,                  // service name
		client.WithRPCTimeout(3*time.Second),              // rpc timeout
		client.WithConnectTimeout(50*time.Millisecond),    // conn timeout
		client.WithFailureRetry(retry.NewFailurePolicy()), // retry
		client.WithResolver(r),                            // resolver
	)
	if err != nil {
		panic(err)
	}
	aiClient = c
}

// OrderQuery 查询订单详情
// orderId 订单ID
// 返回AI格式化后的订单信息
func OrderQuery(ctx context.Context, orderId string) (string, error) {
	req := &ai.OrderQueryReq{
		OrderId: orderId,
	}
	resp, err := aiClient.OrderQuery(ctx, req)
	if err != nil {
		return "", err
	}
	return resp.Response, nil
}

// AutoPlaceOrder 自动下单服务
// userId 用户ID
// request 用户的下单请求描述
// 返回订单ID
func AutoPlaceOrder(ctx context.Context, userId uint32, request string) (string, error) {
	req := &ai.AutoPlaceOrderReq{
		UserId:  userId,
		Request: request,
	}
	resp, err := aiClient.AutoPlaceOrder(ctx, req)
	if err != nil {
		return "", err
	}
	return resp.OrderId, nil
}
