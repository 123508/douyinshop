package client

import (
	"context"
	"github.com/123508/douyinshop/kitex_gen/payment"
	"github.com/123508/douyinshop/kitex_gen/payment/paymentservice"
	"log"

	"github.com/123508/douyinshop/pkg/config"

	"time"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/retry"
	etcd "github.com/kitex-contrib/registry-etcd"
)

var paymentClient paymentservice.Client

func initPaymentRpc() {
	r, err := etcd.NewEtcdResolverWithAuth(config.Conf.EtcdConfig.Endpoints, config.Conf.EtcdConfig.Username, config.Conf.EtcdConfig.Password)
	if err != nil {
		panic(err)
	}

	c, err := paymentservice.NewClient(
		config.Conf.PaymentConfig.ServiceName,             // service name
		client.WithRPCTimeout(3*time.Second),              // rpc timeout
		client.WithConnectTimeout(50*time.Millisecond),    // conn timeout
		client.WithFailureRetry(retry.NewFailurePolicy()), // retry
		client.WithResolver(r),                            // resolver
	)
	if err != nil {
		panic(err)
	}
	paymentClient = c
}

func Charge(ctx context.Context, amount float32, orderId string, userId uint32, payMethod int32) (resp *payment.ChargeResp, err error) {
	req := &payment.ChargeReq{
		Amount:    amount,
		OrderId:   orderId,
		UserId:    userId,
		PayMethod: payMethod,
	}
	resp, err = paymentClient.Charge(ctx, req)
	if err != nil {
		log.Println(err)
		return
	}
	return
}

func Notify(ctx context.Context, orderId string, tracsactionId string) {
	req := &payment.NotifyReq{
		OrderId:       orderId,
		TransactionId: tracsactionId,
	}
	_, err := paymentClient.Notify(ctx, req)
	if err != nil {
		log.Println(err)
		return
	}
	return
}
