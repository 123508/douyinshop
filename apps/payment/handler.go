package main

import (
	"context"
	"fmt"
	payment "github.com/123508/douyinshop/kitex_gen/payment"
	"github.com/123508/douyinshop/pkg/config"
	"github.com/smartwalle/alipay/v3"
	"math/rand"
	"sync/atomic"
	"time"
)

// PaymentServiceImpl implements the last service interface defined in the IDL.
type PaymentServiceImpl struct{}

// Charge implements the PaymentServiceImpl interface.
func (s *PaymentServiceImpl) Charge(ctx context.Context, req *payment.ChargeReq) (resp *payment.ChargeResp, err error) {
	if req.PayMethod == 2 {
		// 初始化支付宝
		var client *alipay.Client
		if client, err = alipay.New(config.Conf.PaymentConfig.Alipay.AppId, config.Conf.PaymentConfig.Alipay.PrivateKey,
			false); err != nil {
			return nil, err
		}
		// 支付链接参数
		var p = alipay.TradePagePay{}
		p.Subject = "抖音商城购物"
		p.OutTradeNo = req.OrderId
		p.TotalAmount = fmt.Sprintf("%.2f", req.Amount)
		p.ProductCode = "FAST_INSTANT_TRADE_PAY"
		p.ReturnURL = ""
		p.NotifyURL = ""
		// 生成支付链接
		orderUrl, _ := client.TradePagePay(p)

		return &payment.ChargeResp{
			TransactionId: "",
			PayUrl:        orderUrl.String(),
		}, nil
	} else {
		// 获取当前日期并格式化为 YYYYMMDD
		currentDay := time.Now().Format("20060102")
		// 生成一个12位的随机数
		var counter uint64
		count := atomic.AddUint64(&counter, 1)
		countStr := fmt.Sprintf("%06d", count%1000000)         // 取计数器的后6位
		randomPart1 := fmt.Sprintf("%06d", rand.Intn(1000000)) // 生成6位随机数
		// 生成一个8位的随机数
		randomPart2 := fmt.Sprintf("%08d", rand.Intn(100000000))
		// 拼接成28位的订单号
		tradeNo := currentDay + countStr + randomPart1 + randomPart2

		return &payment.ChargeResp{
			TransactionId: tradeNo,
			PayUrl:        "",
		}, nil
	}
}

// Notify implements the PaymentServiceImpl interface.
func (s *PaymentServiceImpl) Notify(ctx context.Context, req *payment.NotifyReq) (resp *payment.NotifyResp, err error) {
	return
}
