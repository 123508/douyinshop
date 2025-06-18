package shop

import (
	"context"
	"github.com/123508/douyinshop/kitex_gen/shop"
	"github.com/123508/douyinshop/pkg/errorno"
	"github.com/123508/douyinshop/pkg/infras/client"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/cloudwego/hertz/pkg/app"
)

func GetShopId(ctx context.Context, c *app.RequestContext) {
	userId, ok := ctx.Value("userId").(uint64)
	if !ok {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "userId must be a number",
		})
		return
	}
	req := &shop.GetShopIdReq{UserId: userId}
	resp, err := client.GetShopId(ctx, req)
	if err != nil {
		errorno.DealWithError(err, c)
		return
	}
	c.JSON(consts.StatusOK, utils.H{
		"shop_id": resp.ShopId,
	})
}
