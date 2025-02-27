package shop

import (
	"context"
	"github.com/123508/douyinshop/apps/api/infras/client"
	"github.com/123508/douyinshop/kitex_gen/shop"

	"github.com/cloudwego/hertz/pkg/app"
)

func UpdateShopInfo(ctx context.Context, c *app.RequestContext) {
	userID, _ := c.Get("userId")
	type Shop struct {
		ID          uint32 `json:"shop_id"`
		Name        string `json:"name"`
		Address     string `json:"address"`
		Description string `json:"description"`
		Avatar      string `json:"avatar"`
	}
	shopInfo := &Shop{}
	err := c.Bind(shopInfo)
	if err != nil {
		c.JSON(400, map[string]interface{}{
			"error": "参数错误",
		})
		return
	}
	getShopIdReq := shop.GetShopIdReq{
		UserId: userID.(uint32),
	}
	getShopIdResp, err := client.GetShopId(ctx, &getShopIdReq)
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	if shopInfo.ID != getShopIdResp.ShopId {
		c.JSON(403, map[string]interface{}{
			"error": "没有权限",
		})
		return
	}
	req := shop.UpdateShopInfoReq{
		ShopId:          shopInfo.ID,
		ShopName:        shopInfo.Name,
		ShopAddress:     shopInfo.Address,
		ShopDescription: shopInfo.Description,
		ShopAvatar:      shopInfo.Avatar,
	}
	resp, err := client.UpdateShopInfo(ctx, &req)
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, map[string]interface{}{
		"ok": resp.Res,
	})
}
