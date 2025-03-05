package ai

import (
	"context"
	"github.com/123508/douyinshop/apps/api/infras/client"
	"github.com/123508/douyinshop/pkg/errors"
	"github.com/cloudwego/hertz/pkg/app"
)

type AutoPlaceOrderRequest struct {
	UserId        uint32 `json:"user_id" vd:"$>0"`                // 用户ID
	Request       string `json:"request" vd:"len($)>0"`           // 下单请求描述
	AddressBookId uint32 `json:"address_book_id,omitempty"`      // 可选：地址ID
	PayMethod     int32  `json:"pay_method,omitempty"`           // 可选：支付方式
}

// AutoPlaceOrder 自动下单
// @router /api/ai/order/place [POST]
func AutoPlaceOrder(ctx context.Context, c *app.RequestContext) {
	var req AutoPlaceOrderRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(400, errors.NewValidationError(err))
		return
	}

	// 调用AI服务自动下单
	resp, err := client.AutoPlaceOrder(ctx, &client.AutoPlaceOrderParams{
		UserId:        req.UserId,
		Request:       req.Request,
		AddressBookId: req.AddressBookId,
		PayMethod:     req.PayMethod,
	})
	if err != nil {
		c.JSON(500, errors.NewServiceError(err))
		return
	}

	c.JSON(200, map[string]interface{}{
		"code": 0,
		"msg":  "success",
		"data": map[string]interface{}{
			"order_id":     resp.OrderId,
			"total_amount": resp.TotalAmount,
			"message":      resp.Message,
		},
	})
}