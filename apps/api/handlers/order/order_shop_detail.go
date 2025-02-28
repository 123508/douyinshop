package order

import (
	"context"
	"github.com/123508/douyinshop/apps/api/infras/client"
	"github.com/123508/douyinshop/kitex_gen/order/order_common"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
)

func DetailShop(ctx context.Context, c *app.RequestContext) {

	orderId, err := strconv.Atoi(c.Query("orderId"))
	if err != nil {
		c.JSON(400, map[string]interface{}{
			"error": "orderId 参数错误",
		})
		return
	}

	var orderDetails []order_common.OrderDetail
	if err := c.Bind(&orderDetails); err != nil {
		c.JSON(400, map[string]interface{}{
			"error": "订单详情解析失败",
		})
		return
	}

	orderResp, err := client.ShopDetail(ctx, uint32(orderId), orderDetails)
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"error": "internal server error",
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"order": orderResp,
	})

}
