package shop

import (
	"context"
	"github.com/123508/douyinshop/apps/api/infras/client"
	"github.com/123508/douyinshop/kitex_gen/shop"

	"github.com/cloudwego/hertz/pkg/app"
)

func GetShopId(ctx context.Context, c *app.RequestContext) {
	userID, _ := c.Get("userId")
	req := &shop.GetShopIdReq{UserId: userID.(uint32)}
	resp, err := client.GetShopId(ctx, req)
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, map[string]interface{}{
		"shop_id": resp.ShopId,
	})
}
