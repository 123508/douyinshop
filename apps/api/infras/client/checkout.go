package client

import (
	"context"
	"github.com/123508/douyinshop/kitex_gen/checkout"
	"github.com/123508/douyinshop/kitex_gen/checkout/checkoutservice"

	"github.com/123508/douyinshop/pkg/config"

	"time"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/retry"
	etcd "github.com/kitex-contrib/registry-etcd"
)

var checkoutClient checkoutservice.Client

func initCheckoutRpc() {
	r, err := etcd.NewEtcdResolverWithAuth(config.Conf.EtcdConfig.Endpoints, config.Conf.EtcdConfig.Username, config.Conf.EtcdConfig.Password)
	if err != nil {
		panic(err)
	}

	c, err := checkoutservice.NewClient(
		config.Conf.CheckoutConfig.ServiceName,            // service name
		client.WithRPCTimeout(3*time.Second),              // rpc timeout
		client.WithConnectTimeout(50*time.Millisecond),    // conn timeout
		client.WithFailureRetry(retry.NewFailurePolicy()), // retry
		client.WithResolver(r),                            // resolver
	)
	if err != nil {
		panic(err)
	}
	checkoutClient = c
}
func Checkout(ctx context.Context, req *checkout.CheckoutReq) (bool, error) {
	_, err := checkoutClient.Checkout(ctx, req)
	if err != nil {
		return false, err
	}
	return true, nil
}
