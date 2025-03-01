package order

import (
	"context"
	"github.com/123508/douyinshop/apps/api/infras/client"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
)

func CancelShop(ctx context.Context, c *app.RequestContext) {

	orderId, err := strconv.Atoi(c.Query("orderId"))
	if err != nil {
		c.JSON(400, map[string]interface{}{
			"error": "orderId 参数错误",
		})
		return
	}

	cancelReason := c.Query("cancelReason")
	if cancelReason == "" {
		c.JSON(400, map[string]interface{}{
			"error": "取消原因不能为空",
		})
		return
	}

	resp, err := client.ShopCancel(ctx, uint32(orderId), cancelReason)
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"error": "订单取消失败",
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"message": "订单取消成功",
		"data":    resp,
	})

}
