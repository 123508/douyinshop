package ai

import (
	"context"
	"github.com/123508/douyinshop/apps/api/infras/client"
	"github.com/123508/douyinshop/pkg/errors"
	"github.com/cloudwego/hertz/pkg/app"
)

type OrderQueryRequest struct {
	OrderId string `json:"order_id" query:"order_id" vd:"$!=''"` // 订单ID
}

// OrderQuery 查询订单详情
// @router /api/ai/order/query [GET]
func OrderQuery(ctx context.Context, c *app.RequestContext) {
	var req OrderQueryRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(400, errors.NewValidationError(err))
		return
	}

	// 调用AI服务查询订单
	response, err := client.OrderQuery(ctx, req.OrderId)
	if err != nil {
		c.JSON(500, errors.NewServiceError(err))
		return
	}

	c.JSON(200, map[string]interface{}{
		"code": 0,
		"msg":  "success",
		"data": map[string]interface{}{
			"response": response,
		},
	})
}