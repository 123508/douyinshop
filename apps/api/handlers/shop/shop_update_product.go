package shop

import (
	"context"
	"github.com/123508/douyinshop/apps/api/infras/client"
	"github.com/123508/douyinshop/kitex_gen/product"
	"github.com/123508/douyinshop/kitex_gen/shop"
	"github.com/cloudwego/hertz/pkg/app"
)

func UpdateProductInfo(ctx context.Context, c *app.RequestContext) {
	userID, _ := c.Get("userId")
	type ProductInfo struct {
		ShopID      uint32   `json:"shop_id"`
		ProductID   uint32   `json:"product_id"`
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
	req := &shop.UpdateProductReq{
		ShopId: productInfo.ShopID,
		Product: &product.Product{
			Id:          productInfo.ProductID,
			Name:        productInfo.Name,
			Description: productInfo.Description,
			Picture:     productInfo.Picture,
			Price:       productInfo.Price,
			Categories:  productInfo.Categories,
			Status:      productInfo.Status,
		},
	}
	resp, err := client.UpdateProduct(ctx, req)
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, map[string]interface{}{
		"ok": resp,
	})
}
