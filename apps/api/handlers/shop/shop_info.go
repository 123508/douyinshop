package shop

import (
	"context"
	"github.com/123508/douyinshop/apps/api/infras/client"
	"github.com/123508/douyinshop/kitex_gen/shop"
	"github.com/123508/douyinshop/pkg/errorno"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
)

func GetInfo(ctx context.Context, c *app.RequestContext) {
	shopID, err := strconv.Atoi(c.Param("shop_id"))
	if err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "shop_id参数错误",
		})
		return
	}
	req := &shop.GetShopInfoReq{
		ShopId: uint32(shopID),
	}
	shopInfo, err := client.GetShopInfo(ctx, req)
	if err != nil {
		basicErr := errorno.ParseBasicMessageError(err)

		if basicErr.Raw != nil {
			c.JSON(consts.StatusInternalServerError, utils.H{
				"err": err,
			})
		} else {
			c.JSON(basicErr.Code, utils.H{
				"error": basicErr.Message,
			})
		}

		return
	}
	c.JSON(consts.StatusOK, utils.H{
		"name":        shopInfo.ShopName,
		"address":     shopInfo.ShopAddress,
		"description": shopInfo.ShopDescription,
		"avatar":      shopInfo.ShopAvatar,
	})
}
