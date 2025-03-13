package shop

import (
	"context"
	"github.com/123508/douyinshop/apps/api/infras/client"
	"github.com/123508/douyinshop/kitex_gen/product"
	"github.com/123508/douyinshop/kitex_gen/shop"
	"github.com/123508/douyinshop/pkg/errorno"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/cloudwego/hertz/pkg/app"
)

func Add(ctx context.Context, c *app.RequestContext) {
	value, exists := c.Get("userId")
	userId, ok := value.(uint32)
	if !exists || !ok {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "userId must be a number",
		})
		return
	}
	type ProductInfo struct {
		ShopID      uint32   `json:"shop_id"`
		Name        string   `json:"product_name"`
		Description string   `json:"product_description"`
		Picture     string   `json:"product_picture"`
		Price       float32  `json:"price"`
		Categories  []string `json:"categories"`
		Status      bool     `json:"status"`
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
	req := &shop.AddProductReq{
		ShopId: productInfo.ShopID,
		Product: &product.Product{
			Name:        productInfo.Name,
			Description: productInfo.Description,
			Picture:     productInfo.Picture,
			Price:       productInfo.Price,
			Categories:  productInfo.Categories,
			Status:      productInfo.Status,
		},
	}
	resp, err := client.AddProduct(ctx, req)
	if err != nil {
		errorno.DealWithError(err, c)
		return
	}
	c.JSON(consts.StatusOK, utils.H{
		"ok": resp.ProductId != 0,
	})
}
