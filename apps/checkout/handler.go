package main

import (
	"context"
	"github.com/123508/douyinshop/kitex_gen/cart"
	"github.com/123508/douyinshop/kitex_gen/cart/cartservice"
	_ "github.com/123508/douyinshop/kitex_gen/cart/cartservice"
	_ "github.com/123508/douyinshop/kitex_gen/order/order_common"
	"github.com/123508/douyinshop/kitex_gen/order/userOrder"
	"github.com/123508/douyinshop/kitex_gen/order/userOrder/orderuserservice"
	"github.com/123508/douyinshop/kitex_gen/payment"
	"github.com/123508/douyinshop/kitex_gen/payment/paymentservice"
	"github.com/123508/douyinshop/kitex_gen/product"
	"github.com/123508/douyinshop/kitex_gen/product/productcatalogservice"
	"github.com/123508/douyinshop/pkg/errorno"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/klog"
	_ "go.opentelemetry.io/otel"
	_ "go.opentelemetry.io/otel/propagation"
	"strconv"
	"sync"

	checkout "github.com/123508/douyinshop/kitex_gen/checkout"
)

// CheckoutServiceImpl implements the last service interface defined in the IDL.
type CheckoutServiceImpl struct {
	ctx context.Context
}

var (
	CartClient    cartservice.Client
	ProductClient productcatalogservice.Client
	PaymentClient paymentservice.Client
	OrderClient   orderuserservice.Client
	once          sync.Once
	err           error
	registryAddr  string
	serviceName   string
	commonSuite   client.Option
)

var GetCartError = &errorno.BasicMessageError{Code: 404, Message: "获取购物车错误"}

var CartNotExistError = &errorno.BasicMessageError{Code: 404, Message: "购物车不存在"}

var SubmitError = &errorno.BasicMessageError{Code: 500, Message: "提交错误"}

var ChargeError = &errorno.BasicMessageError{Code: 417, Message: "支付异常"}

// Checkout implements the CheckoutServiceImpl interface.
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
	var orderId uint32
	if orderResult != nil || orderResult.OrderId != 0 {
		orderId = orderResult.OrderId
	}
	payReq := &payment.ChargeReq{
		UserId:  req.UserId,
		OrderId: strconv.Itoa(int(orderId)),
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
		OrderId:       strconv.Itoa(int(orderId)),
		TransactionId: paymentResult.TransactionId,
	}
	return

}
