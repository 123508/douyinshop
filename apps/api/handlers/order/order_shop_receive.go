package order

import (
	"context"
	"github.com/123508/douyinshop/apps/api/infras/client"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
)

func Receive(ctx context.Context, c *app.RequestContext) {

	orderId, err := strconv.Atoi(c.Query("orderId"))
	if err != nil {
		c.JSON(400, map[string]interface{}{
			"error": "orderId 参数错误",
		})
		return
	}

	resp, err := client.ShopReceive(ctx, uint32(orderId))
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"error": "订单确认收货失败",
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"message": "订单确认收货成功",
		"data":    resp,
	})

}
