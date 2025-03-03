package shop

import (
	"context"

	"github.com/123508/douyinshop/apps/api/infras/client"
	"github.com/123508/douyinshop/kitex_gen/shop"

	"github.com/cloudwego/hertz/pkg/app"
)

func Register(ctx context.Context, c *app.RequestContext) {
	userID, _ := c.Get("userId")
	type Shop struct {
		Name        string `json:"name"`
		Address     string `json:"address"`
		Description string `json:"description"`
		Avatar      string `json:"avatar"`
	}
	var shopInfo Shop
	err := c.Bind(&shopInfo)
	if err != nil {
		c.JSON(400, map[string]interface{}{
			"error": "参数错误",
		})
		return
	}
	req := shop.RegisterShopReq{
		UserId:          userID.(uint32),
		ShopName:        shopInfo.Name,
		ShopAddress:     shopInfo.Address,
		ShopDescription: shopInfo.Description,
		ShopAvatar:      shopInfo.Avatar,
	}
	resp, err := client.RegisterShop(ctx, &req)
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, map[string]interface{}{
		"ok": resp.ShopId != 0,
	})
}
