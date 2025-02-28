package order

import (
	"context"
	"github.com/123508/douyinshop/apps/api/infras/client"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
)

func Cancel(ctx context.Context, c *app.RequestContext) {

	orderId, err := strconv.Atoi(c.Query("orderId"))
	if err != nil {
		c.JSON(400, map[string]interface{}{
			"error": "orderId 参数错误",
		})
		return
	}

	CancelReason := c.Query("cancelReason")
	if CancelReason == "" {
		c.JSON(400, map[string]interface{}{
			"error": "cancelReason 参数不能为空",
		})
		return
	}

	resp, err := client.UserCancel(ctx, uint32(orderId), CancelReason)
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"error": "internal server error",
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"message": "订单已取消",
		"result":  &resp,
	})

}
