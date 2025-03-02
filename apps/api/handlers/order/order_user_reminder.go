package order

import (
	"context"
	"github.com/123508/douyinshop/apps/api/infras/client"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
)

func Reminder(ctx context.Context, c *app.RequestContext) {

	userId, err := strconv.Atoi(c.Query("userId"))
	if err != nil {
		c.JSON(400, map[string]interface{}{
			"error": "userId 参数错误",
		})
		return
	}

	orderId, err := strconv.Atoi(c.Query("orderId"))
	if err != nil {
		c.JSON(400, map[string]interface{}{
			"error": "orderId 参数错误",
		})
		return
	}

	resp, err := client.UserReminder(ctx, uint32(userId), uint32(orderId))
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"error": "internal server error",
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"message": "已成功提醒商家发货",
		"result":  &resp,
	})

}
