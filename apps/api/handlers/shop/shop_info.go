package shop

import (
	"context"
	"github.com/123508/douyinshop/apps/api/infras/client"
	"github.com/123508/douyinshop/kitex_gen/shop"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
)

func GetInfo(ctx context.Context, c *app.RequestContext) {
	shopID, err := strconv.Atoi(c.Param("shop_id"))
	if err != nil {
		c.JSON(400, map[string]interface{}{
			"error": "shop_id参数错误",
		})
		return
	}
	req := &shop.GetShopInfoReq{
		ShopId: uint32(shopID),
	}
	shopInfo, err := client.GetShopInfo(ctx, req)
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, map[string]interface{}{
		"name":        shopInfo.ShopName,
		"address":     shopInfo.ShopAddress,
		"description": shopInfo.ShopDescription,
		"avatar":      shopInfo.ShopAvatar,
	})
}
