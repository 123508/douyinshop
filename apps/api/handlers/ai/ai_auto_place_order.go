package ai

import (
	"context"
	"github.com/123508/douyinshop/apps/api/infras/client"
	"github.com/123508/douyinshop/pkg/errors"
	"github.com/cloudwego/hertz/pkg/app"
)

type AutoPlaceOrderRequest struct {
	UserId  uint32 `json:"user_id" vd:"$>0"`      // 用户ID
	Request string `json:"request" vd:"len($)>0"`  // 下单请求描述
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
	orderId, err := client.AutoPlaceOrder(ctx, req.UserId, req.Request)
	if err != nil {
		c.JSON(500, errors.NewServiceError(err))
		return
	}

	c.JSON(200, map[string]interface{}{
		"code": 0,
		"msg":  "success",
		"data": map[string]interface{}{
			"order_id": orderId,
		},
	})
}