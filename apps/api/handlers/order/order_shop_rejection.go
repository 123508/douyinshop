package order

import (
	"context"
	"github.com/123508/douyinshop/apps/api/infras/client"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
)

func Rejection(ctx context.Context, c *app.RequestContext) {

	orderId, err := strconv.Atoi(c.Query("orderId"))
	if err != nil {
		c.JSON(400, map[string]interface{}{
			"error": "orderId 参数错误",
		})
		return
	}

	rejectionReason := c.Query("rejectionReason")
	if rejectionReason == "" {
		c.JSON(400, map[string]interface{}{
			"error": "拒绝原因不能为空",
		})
		return
	}

	resp, err := client.ShopRejection(ctx, uint32(orderId), rejectionReason)
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"error": "订单拒绝失败",
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"message": "订单拒绝成功",
		"data":    resp,
	})

}
