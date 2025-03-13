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

func Delete(ctx context.Context, c *app.RequestContext) {
	value, exists := c.Get("userId")
	userId, ok := value.(uint32)
	if !exists || !ok {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "userId must be a number",
		})
		return
	}
	type ProductInfo struct {
		ShopID    uint32 `json:"shop_id"`
		ProductID uint32 `json:"product_id"`
	}
	productInfo := &ProductInfo{}
	err := c.Bind(productInfo)
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
	if productInfo.ShopID != getShopIdResp.ShopId {
		c.JSON(consts.StatusForbidden, utils.H{
			"error": "没有权限",
		})
		return
	}
	req := &shop.DeleteProductReq{
		ShopId:    productInfo.ShopID,
		ProductId: productInfo.ProductID,
	}
	resp, err := client.DeleteProduct(ctx, req)
	if err != nil {
		errorno.DealWithError(err, c)
		return
	}
	c.JSON(consts.StatusOK, utils.H{
		"ok": resp.Res,
	})
}
