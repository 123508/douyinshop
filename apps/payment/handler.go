package main

import (
	"context"
	"fmt"
	payment "github.com/123508/douyinshop/kitex_gen/payment"
	"github.com/123508/douyinshop/pkg/config"
	"github.com/123508/douyinshop/pkg/errors"
	"github.com/123508/douyinshop/pkg/models"
	"github.com/smartwalle/alipay/v3"
	"log"
	"math/rand"
	"sync/atomic"
	"time"
)

// PaymentServiceImpl implements the last service interface defined in the IDL.
type PaymentServiceImpl struct{}

// Charge implements the PaymentServiceImpl interface.
func (s *PaymentServiceImpl) Charge(ctx context.Context, req *payment.ChargeReq) (resp *payment.ChargeResp, err error) {
	var order models.Order
	if err = database.Model(&models.Order{}).Where("user_id = ? and number = ?", req.UserId, req.OrderId).First(&order).Error; err != nil {
		return nil, &errors.BasicMessageError{Message: "没有该购物记录"}
	}
	if req.PayMethod == 1 {
		return nil, fmt.Errorf("暂时不支持微信支付")
	} else if req.PayMethod == 2 {
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

		//修改数据库中对应记录
		// Order表的更改:PayStatus和TransactionId,获取其主键
		updateFields := map[string]interface{}{
			"pay_status":     1,
			"transaction_id": tradeNo,
		}
		if err := database.Model(&order).Updates(updateFields).Error; err != nil {
			return nil, err
		}
		Id := order.ID
		// OrderDetail表的更改:Status,获取其主键
		result := database.Model(&models.OrderDetail{}).Where("order_id = ?", Id).Update("status", 1)
		if result.Error != nil {
			log.Println(result.Error)
		}
		var orderDetail models.OrderDetail
		result = database.Model(&models.OrderDetail{}).Where("order_id = ?", Id).First(&orderDetail)
		Id = orderDetail.ID

		// OrderStatusLog表的更改:改一个，加一个
		result = database.Model(&models.OrderStatusLog{}).Where("order_detail_id = ? AND status = ?", Id, 0).Update("end_time", time.Now())
		if result.Error != nil {
			log.Println(result.Error)
		}
		result = database.Create(&models.OrderStatusLog{
			OrderDetailId: int(Id),
			Status:        1,
			StartTime:     time.Now(),
		})
		if result.Error != nil {
			log.Println(result.Error)
		}

		return &payment.ChargeResp{
			TransactionId: tradeNo,
			PayUrl:        "",
		}, nil
	}
}

// Notify implements the PaymentServiceImpl interface.
func (s *PaymentServiceImpl) Notify(ctx context.Context, req *payment.NotifyReq) (resp *payment.NotifyResp, err error) {
	orderId := req.OrderId
	transactionId := req.TransactionId

	// Order表的更改:PayStatus和TransactionId,获取其主键
	result := database.Model(&models.Order{}).Where("number = ?", orderId).First(&models.Order{}).
		Updates(map[string]interface{}{"pay_status": 1, "transaction_id": transactionId})
	if result.Error != nil {
		log.Println(result.Error)
	}
	var order models.Order
	result = database.Model(&models.Order{}).Where("number = ?", orderId).First(&order)
	Id := order.ID

	// OrderDetail表的更改:Status,获取其主键
	result = database.Model(&models.OrderDetail{}).Where("order_id = ?", Id).Update("status", 1)
	if result.Error != nil {
		log.Println(result.Error)
	}
	var orderDetail models.OrderDetail
	result = database.Model(&models.OrderDetail{}).Where("order_id = ?", Id).First(&orderDetail)
	Id = orderDetail.ID

	// OrderStatusLog表的更改:改一个，加一个
	result = database.Model(&models.OrderStatusLog{}).Where("order_detail_id = ? AND status = ?", Id, 0).Update("end_time", time.Now())
	if result.Error != nil {
		log.Println(result.Error)
	}
	result = database.Create(&models.OrderStatusLog{
		OrderDetailId: int(Id),
		Status:        1,
		StartTime:     time.Now(),
	})
	if result.Error != nil {
		log.Println(result.Error)
	}
	return &payment.NotifyResp{}, nil
}
