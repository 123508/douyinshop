package main

import (
	"context"
	"github.com/123508/douyinshop/kitex_gen/cart"
	_ "github.com/123508/douyinshop/kitex_gen/cart/cartservice"
	"github.com/123508/douyinshop/kitex_gen/checkout"
	_ "github.com/123508/douyinshop/kitex_gen/order/order_common"
	"github.com/123508/douyinshop/kitex_gen/order/userOrder"
	"github.com/123508/douyinshop/kitex_gen/payment"
	"github.com/123508/douyinshop/kitex_gen/product"
	"github.com/cloudwego/kitex/pkg/klog"
	_ "go.opentelemetry.io/otel"
	_ "go.opentelemetry.io/otel/propagation"
)

// CheckoutServiceImpl implements the last service interface defined in the IDL.
type CheckoutServiceImpl struct {
	ctx context.Context
}

// Checkout implements the CheckoutServiceImpl interface.
// 结算接口
func (s *CheckoutServiceImpl) Checkout(ctx context.Context, req *checkout.CheckoutReq) (resp *checkout.CheckoutResp, err error) {
	//get cart
	cartResult, err := CartClient.GetCart(s.ctx, &cart.GetCartReq{UserId: req.UserId})
	if err != nil {
		klog.Error("获取购物车错误:", err)
		return nil, GetCartError
	}
	if cartResult == nil || cartResult.Cart == nil || len(cartResult.Cart.Items) == 0 {
		klog.Error("购物车不存在")
		return nil, CartNotExistError
	}
	var (
		//oi    []*order_common.Order
		total float32
	)
	for _, cartItem := range cartResult.Cart.Items {
		productResp, resultErr := ProductClient.GetProduct(s.ctx, &product.GetProductReq{Id: cartItem.ProductId})
		if resultErr != nil {
			klog.Error(resultErr)
			return nil, resultErr
		}
		if productResp.Product == nil {
			continue
		}
		p := productResp.Product
		cost := p.Price * float32(cartItem.Quantity)
		total += cost
		/*oi = append(oi, &order_common.Order{
		}*/
	}
	//create order
	orderReq := &userOrder.OrderSubmitReq{
		UserId:    req.UserId,
		PayMethod: 1, //不知道如何传入
		Remark:    "0",
		Amount:    total,
	}

	if req.Address != nil {
		orderReq.AddressBookId = 1 //不知道如何传入
	}
	orderResult, err := OrderClient.Submit(s.ctx, orderReq)
	if err != nil {
		klog.Error("提交错误:", err)
		return nil, SubmitError
	}
	klog.Info("orderResult", orderResult)
	// empty cart
	emptyResult, err := CartClient.EmptyCart(s.ctx, &cart.EmptyCartReq{UserId: req.UserId})
	if err != nil {
		klog.Error("购物车不存在:", err)
		return nil, CartNotExistError
	}
	klog.Info(emptyResult)
	// charge
	var orderId uint64
	if orderResult != nil && orderResult.OrderId != 0 {
		orderId = orderResult.OrderId
	}
	payReq := &payment.ChargeReq{
		UserId:  req.UserId,
		OrderId: orderId,
		Amount:  total,
		CreditCard: &payment.CreditCardInfo{
			CreditCardNumber:          req.CreditCard.CreditCardNumber,
			CreditCardExpirationYear:  req.CreditCard.CreditCardExpirationYear,
			CreditCardExpirationMonth: req.CreditCard.CreditCardExpirationMonth,
			CreditCardCvv:             req.CreditCard.CreditCardCvv,
		},
	}
	paymentResult, err := PaymentClient.Charge(s.ctx, payReq)
	if err != nil {
		klog.Error("支付异常:", err)
		return nil, ChargeError
	}

	// otel inject

	klog.Info(paymentResult)
	// change order state
	klog.Info(orderResult)

	resp = &checkout.CheckoutResp{
		OrderId:       orderId,
		TransactionId: paymentResult.TransactionId,
	}
	return

}
