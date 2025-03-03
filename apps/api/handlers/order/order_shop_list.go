package order

import (
	"context"
	"github.com/123508/douyinshop/apps/api/infras/client"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
)

func List(ctx context.Context, c *app.RequestContext) {

	shopId, err := strconv.Atoi(c.Query("shopId"))
	if err != nil {
		c.JSON(400, map[string]interface{}{
			"error": "shopId 参数错误",
		})
		return
	}

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		c.JSON(400, map[string]interface{}{
			"error": "page 参数错误",
		})
		return
	}

	pageSize, err := strconv.Atoi(c.Query("pageSize"))
	if err != nil {
		c.JSON(400, map[string]interface{}{
			"error": "pageSize 参数错误",
		})
		return
	}

	orderResp, err := client.GetOrderList(ctx, uint32(shopId), uint32(page), uint32(pageSize))
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"error": "internal server error",
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"orders": orderResp,
	})
}
