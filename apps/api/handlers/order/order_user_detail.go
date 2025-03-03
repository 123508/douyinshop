package order

import (
	"context"
	"github.com/123508/douyinshop/apps/api/infras/client"
	"github.com/123508/douyinshop/pkg/models"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
)

func Detail(ctx context.Context, c *app.RequestContext) {

	orderId, err := strconv.Atoi(c.Query("orderId"))
	if err != nil {
		c.JSON(400, map[string]interface{}{
			"error": "orderId参数错误",
		})
		return
	}
	var orderDetails []models.OrderDetail

	orderResp, err := client.UserDetail(ctx, uint32(orderId), orderDetails)
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
