package shop

import (
	"context"
	"github.com/123508/douyinshop/apps/api/infras/client"
	"github.com/123508/douyinshop/kitex_gen/shop"

	"github.com/cloudwego/hertz/pkg/app"
)

func Delete(ctx context.Context, c *app.RequestContext) {
	userID, _ := c.Get("userId")
	type ProductInfo struct {
		ShopID    uint32 `json:"shop_id"`
		ProductID uint32 `json:"product_id"`
	}
	productInfo := &ProductInfo{}
	err := c.Bind(productInfo)
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
	if productInfo.ShopID != getShopIdResp.ShopId {
		c.JSON(403, map[string]interface{}{
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
		c.JSON(500, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, map[string]interface{}{
		"ok": resp.Res,
	})
}
