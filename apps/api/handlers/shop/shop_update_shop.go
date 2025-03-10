package shop

import (
	"context"
	"github.com/123508/douyinshop/apps/api/infras/client"
	"github.com/123508/douyinshop/kitex_gen/shop"
	"github.com/123508/douyinshop/pkg/errorno"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/cloudwego/hertz/pkg/app"
)

func UpdateShopInfo(ctx context.Context, c *app.RequestContext) {
	value, exists := c.Get("userId")
	userId, ok := value.(uint32)
	if !exists || !ok {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "userId must be a number",
		})
		return
	}
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
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "参数错误",
		})
		return
	}
	getShopIdReq := shop.GetShopIdReq{
		UserId: userId,
	}
	getShopIdResp, err := client.GetShopId(ctx, &getShopIdReq)
	if err != nil {
		errorno.DealWithError(err, c)
		return
	}
	if shopInfo.ID != getShopIdResp.ShopId {
		c.JSON(consts.StatusForbidden, utils.H{
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
		errorno.DealWithError(err, c)
		return
	}
	c.JSON(consts.StatusOK, utils.H{
		"ok": resp.Res,
	})
}
