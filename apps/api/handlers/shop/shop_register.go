package shop

import (
	"context"
	"github.com/123508/douyinshop/pkg/errorno"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/123508/douyinshop/apps/api/infras/client"
	"github.com/123508/douyinshop/kitex_gen/shop"

	"github.com/cloudwego/hertz/pkg/app"
)

func Register(ctx context.Context, c *app.RequestContext) {
	userId, ok := ctx.Value("userId").(uint32)
	if !ok {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "userId must be a number",
		})
		return
	}
	type Shop struct {
		Name        string `json:"name"`
		Address     string `json:"address"`
		Description string `json:"description"`
		Avatar      string `json:"avatar"`
	}
	var shopInfo Shop
	err := c.Bind(&shopInfo)
	if err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "参数错误",
		})
		return
	}
	req := shop.RegisterShopReq{
		UserId:          userId,
		ShopName:        shopInfo.Name,
		ShopAddress:     shopInfo.Address,
		ShopDescription: shopInfo.Description,
		ShopAvatar:      shopInfo.Avatar,
	}
	resp, err := client.RegisterShop(ctx, &req)
	if err != nil {
		errorno.DealWithError(err, c)
		return
	}
	c.JSON(consts.StatusOK, utils.H{
		"ok": resp.ShopId != 0,
	})
}
